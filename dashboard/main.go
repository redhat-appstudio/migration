package main

import (
	"encoding/csv"
	"os"
	"strings"

	"context"
	"flag"
	"path/filepath"

	"golang.org/x/sync/errgroup"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/redhat-appstudio/migration/dashboard/pkg/api"

	applicationapiv1alpha1 "github.com/redhat-appstudio/application-api/api/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	integv1alpha1 "github.com/redhat-appstudio/integration-service/api/v1alpha1"
	releasev1alpha1 "github.com/redhat-appstudio/release-service/api/v1alpha1"

	// "k8s.io/apimachinery/pkg/runtime/serializer"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
)

var (
	scheme = runtime.NewScheme()
)

type AppConfigMap map[string]*api.AppConfig
type NsToAppConfigs map[string][]AppConfigMap

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(releasev1alpha1.AddToScheme(scheme))
	utilruntime.Must(applicationapiv1alpha1.AddToScheme(scheme))
	utilruntime.Must(integv1alpha1.AddToScheme(scheme))

	//+kubebuilder:scaffold:scheme
}

func main() {
	restConfig := loadRestConfig()
	cl, err := runtimeclient.New(restConfig, runtimeclient.Options{
		Scheme: scheme,
	})
	if err != nil {
		panic(err)
	}

	ns, err := GetNamespaces(cl, context.Background())
	if err != nil {
		panic(err)
	}

	appConfigs, err := LoadAppConfigs(cl, ns, context.Background())
	if err != nil {
		panic(err)
	}

	AppConfigsToCSV(appConfigs)
}

func loadRestConfig() *rest.Config {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String(
			"kubeconfig",
			filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file",
		)
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	return config
}

func GetNamespaces(cl runtimeclient.Client, ctx context.Context) (*corev1.NamespaceList, error) {
	nsList := &corev1.NamespaceList{}
	err := cl.List(ctx, nsList)

	if err != nil {
		return nil, err
	}
	return nsList, nil
}

func LoadAppConfigs(cl runtimeclient.Client, namespaces *corev1.NamespaceList, ctx context.Context) (NsToAppConfigs, error) {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(3)
	nsToAppConfig := make(NsToAppConfigs)
	for _, ns := range namespaces.Items {

		if !strings.HasSuffix(ns.Name, "-tenant") {
			continue
		}

		nsName := ns.Name

		g.Go(func() error {
			appConfig, err := LoadAppConfigMap(cl, nsName, ctx)
			if err != nil {
				return err
			}
			nsToAppConfig[nsName] = append(nsToAppConfig[nsName], appConfig)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return nsToAppConfig, nil
}

func LoadAppConfigMap(cl runtimeclient.Client, namespace string, ctx context.Context) (AppConfigMap, error) {
	inNamespace := runtimeclient.InNamespace(namespace)
	g, ctx := errgroup.WithContext(ctx)

	apps := &applicationapiv1alpha1.ApplicationList{}
	g.Go(func() error {
		return cl.List(ctx, apps, inNamespace)
	})

	components := &applicationapiv1alpha1.ComponentList{}
	g.Go(func() error {
		return cl.List(ctx, components, inNamespace)
	})

	integTests := &integv1alpha1.IntegrationTestScenarioList{}
	g.Go(func() error {
		return cl.List(ctx, integTests, inNamespace)
	})

	releasePlans := &releasev1alpha1.ReleasePlanList{}
	g.Go(func() error {
		return cl.List(ctx, releasePlans, inNamespace)
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	appConfigMap := make(AppConfigMap)

	if len(apps.Items) == 0 {
		return appConfigMap, nil
	}

	for _, app := range apps.Items {
		appConfigMap[app.Name] = &api.AppConfig{Application: app}
	}

	for _, comp := range components.Items {
		appConfigMap[comp.Spec.Application].AddComponent(comp)
	}

	for _, test := range integTests.Items {
		appConfigMap[test.Spec.Application].AddTest(test)
	}

	for _, rp := range releasePlans.Items {
		appConfigMap[rp.Spec.Application].AddReleasePlan(rp)
	}

	return appConfigMap, nil
}

func AppConfigsToCSV(nsToAppconfigMaps NsToAppConfigs) error {
	f, err := os.Create("result.csv")
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	err = w.Write([]string{
		"Workspace", "AppName", "ComponentName", "IntegrationTests", "ReleasePlans",
	})

	if err != nil {
		return err
	}

	for ns, appConfigMaps := range nsToAppconfigMaps {
		for _, appConfigMap := range appConfigMaps {
			for appName, appConfig := range appConfigMap {
				for _, comp := range appConfig.Components {
					err := w.Write([]string{
						ns,
						appName,
						comp.Name,
						appConfig.GetTestNames(),
						appConfig.GetReleasePlanNames(),
					})
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
