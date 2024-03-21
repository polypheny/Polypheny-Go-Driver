<p align="center">
    <a href="https://polypheny.org/">
        <picture><source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/polypheny/Admin/master/Logo/logo-white-text_cropped.png">
            <img width='50%' alt="Light: 'Resume application project app icon' Dark: 'Resume application project app icon'" src="https://raw.githubusercontent.com/polypheny/Admin/master/Logo/logo-transparent_cropped.png">
        </picture>
    </a>    
</p> 


# Polypheny-DB Go Driver WIP


A Polypheny-DB Driver for the Go programming language.


## Installation

To use this driver in a Go project, simply import the package.

`import "github.com/polypheny/Polypheny-Go-Driver"`.

However, currently the proto files are not hosted at [Polypheny](https://polypheny.com). Instead, they are temporarily included in this repo without being complied. Therefore, before build a Go project, we need to compile these proto files and run `go mod edit -replace polypheny.com/protos=ProtoPath`. 

An in-depth and more detailed documentation can be found [here](https://docs.polypheny.com/en/latest/drivers/go/overview).



## Roadmap
See the [open issues](https://github.com/polypheny/Polypheny-DB/labels/A-golang) for a list of proposed features (and known issues).


## Contributing
We highly welcome your contributions to the _Polypheny Go Driver_. If you would like to contribute, please fork the repository and submit your changes as a pull request. Please consult our [Admin Repository](https://github.com/polypheny/Admin) and our [Website](https://polypheny.org) for guidelines and additional information.

Please note that we have a [code of conduct](https://github.com/polypheny/Admin/blob/master/CODE_OF_CONDUCT.md). Please follow it in all your interactions with the project. 




## License
The Apache 2.0 License
