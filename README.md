# SHU-autoSelection

## Overview
An auto-course-selection program written in Golang.

## Getting started

```sh
git clone https://github.com/Chasing1020/SHU-autoSelection.git
cd SHU-autoSelection
go mod tidy
go get -t .

# *** Then modify the configuration file: `config.yaml` ***

go run main.go # enjoy it
```

And if you want to use a cronjob, you should use `crontab -e` command and add
```
# take Fri Jan 7 20:30:00 CST 2022 as an example
30 20 7 1 * go run /[PATH_TO_CUR]/main.go 2>&1 >> ~/selection.log
```


## Bugs

Have found any bugs or suggestions? 
Please visit the [issue tracker](https://github.com/Chasing1020/SHU-autoSelection/issues).

I'm glad if you have any feedback or give pull requests to this project.


## Disclaimer

This project is only available for free academic discussions.

License [CC BY-NC-SA 4.0](https://creativecommons.org/licenses/by-nc-sa/4.0/deed.en)

## Dependencies

> github.com/gocolly/colly
> 
> gopkg.in/yaml.v2

