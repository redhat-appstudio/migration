# Create app script

Script to deploy several applications on top of Stonesoup. It uses the name of the repo as the name of the application and component, and the repo URL as the git source for the component to be deployed from.

## Usage

Create a file with one Github repo URL per line like the `input-file.txt` in this directory. 

> **IMPORTANT:** The last line of the file has to be an empty line for the script to work.

To use the script, execute it passing as parameter the path to your file.

```bash
./create-app.sh <path-to-your-file>
```