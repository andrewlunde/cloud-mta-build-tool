
### How to get help on the tools commands

| Command | Usage &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;       | Description                                                    
| ------  | --------       |  ----------                                                
| help    | `mbt -h`     | Prints all the available commands.                           
| help    | `mbt [command] --help` or<br> `mbt [command] -h`    | Prints detailed information about the specified command.|

&nbsp;
### How to find out the installed tool version

| Command | Usage &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;       | Description                                                    
| ------  | --------       |  ----------                                                
| version | `mbt -v`     | Prints the current Cloud MTA Build Tool version.                                        <br>

&nbsp;
### How to build an MTA archive from the project sources

#### Prerequisites
* `GNU Make 4.2.1` is installed in your build environment. 
* Module build tools are installed in your build environment.

For more information, see the corresponding [`Download` and `Installation` sections](download.md).

#### Quick start example

```go
// Generates the `Makefile.mta` file.
mbt init 

// Executes the MTA project build for Cloud Foundry target environment.
make -f Makefile.mta p=cf

```

#### Cloud MTA Build Tool commands

| Command | Usage &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;       | Description                                                    
| ------  | --------       |  ----------                                                
| init    | `mbt init <flags>`     | Generates the `Makefile.mta` file according to the MTA descriptor (`mta.yaml` file). <br> The `make` command uses the generated `Makefile.mta` file to package the MTA project. <br> Use the `mbt init` command with the following flags:<br><ul><li>`-s (--source)` is the path to the MTA project; the current path is set as the default.<br> Example: `mbt init -s C:/TestProject` <li>`-t (--target)` is the path to the generated `Makefile` folder; the current path is set as the default. <br> Example: `mbt init -t C:/TestFolder`   

<br>
#### Make commands

The `Makefile.mta` file that is generated by the `mbt init` command is the actual project "builder". It provides the verbose build manifest, which can be changed according to the project needs. It is responsible for:

- Building each of the modules in the MTA project.
- Invoking the MBT commands in the correct order.
- Archiving the MTA project.<br><br>
Use the `make` command to package the MTA project with the following parameters:



| Parameter        | Type | Mandatory&nbsp;/<br>Optional        | Description&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                 | Examples    
| -----------  | ------ | -------       |  ----------                              |  --------------------------------------
| `-f <path to Makefile.mta>`    | string     | Mandatory  | The path to the Makefile.mta file that contains build configurations                             | `make -f Makefile.mta p=cf`
| `p`  | string     | Mandatory     | The name of the deployment platform. <br>The supported deployment platforms are: <ul><li>`cf` for SAP Cloud Platform Cloud Foundry environment  <li>`neo` for the SAP Cloud Platform Neo environment <li>`xsa` for the SAP HANA XS advanced model                                     |`make -f Makefile.mta p=cf`
| `t`    | string     | Optional  | The folder for the generated `MTAR` file. The default value is the current folder. If this parameter is not provided, the `MTAR` file is saved in the `mta_archives` subfolder of the current folder. If the parameter is provided, the `MTAR` file is saved in the root of the folder provided by the argument.                              | `make -f Makefile.mta p=cf t=C:\temp`
| `mtar`    | string     |   Optional  | The file name of the generated archive file. If this parameter is omitted, the file name is created according to the following naming convention: <br><br> `<mta_application_ID>_<mta_application_version>.mtar` <br><br> If the parameter is provided, but does not include an extension, the `.mtar` extension is added. | `make -f Makefile.mta p=cf mtar=myMta`<br><br> `make -f Makefile.mta p=cf mtar=myMta.mtar`
| `strict`    | Boolean     | Optional    | The default value is `true`. If set to `true`, the duplicated fields and fields that are not defined in the `mta.yaml` schema are reported as errors. If set to `false`, they are reported as warnings. | `make -f Makefile.mta p=cf strict=false`

&nbsp;
### How to build an MTA archive from the modules' build artifacts 

| Command | Usage &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;       | Description                                                    
| ------  | --------       |  ----------                                                
| assemble    | `mbt assemble`     | Creates an MTA archive `.mtar` file from the module build artifacts according to the MTA deployment descriptor (`mtad.yaml` file). Runs the command in the directory where the `mtad.yaml` file is located. <br>**Note:** Make sure the path property of each module's `mtad.yaml` file points to the module's build artifacts that you want to package into the target MTA archive. 