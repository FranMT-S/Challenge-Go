# Challenge - Indexer Enron Mail
This project works in conjunction with two other projects. the [Fronted](https://github.com/FranMT-S/fronted) which is responsible for reading the files and the [Backend](https://github.com/FranMT-S/chi-zinc-server)  who provides the query endpoints to the database

 
This project is responsible for indexing emails that come from the [Enron Corp](http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz)


# Enviroment Variable
in the path `src/constanst/env.const.go` there is the environment variable INDEX which is the name of the index used in the database make sure it is the same in the  [Backend](https://github.com/FranMT-S/chi-zinc-server) 

# Run

Execute Command:
``` Go run main.go```

