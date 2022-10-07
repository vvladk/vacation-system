.SILENT:
dbName = ./storage/vs.sqlite
dbInit = ./storage/init

.PHONY: help #  - Generate list of targets with descriptions                                                                
help:    
	@grep '^.PHONY: .* #' ./Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/\1 \2/' | expand -t20

.PHONY: migrate #  - Initialaze the DB

migrate:
	rm $(dbName)  2> /dev/null || true
	$(eval sqlScripts=$(sort $(dbInit)/*.sql))

	for file in $(sqlScripts) ; do \
		sqlite3 $(dbName) < $${file} ; \
    	echo $${file} ; \
    done
.PHONY: gen #  - Create a dunny files controller and model 
gen:
	mkdir ./controller/c$(name)
	# mkdir ./model/m$(name)

	echo "package c$(name)" > ./controller/c$(name)/c$(name).go
	echo "// GET /c$(name)s -  c$(name).index	display a list of all items " >> ./controller/c$(name)/c$(name).go
	echo "// GET /c$(name)/new -  c$(name).new	return an HTML form for creating a new item " >> ./controller/c$(name)/c$(name).go
	echo "// POST /c$(name) -  c$(name).create	create a new item" >> ./controller/c$(name)/c$(name).go
	echo "// GET /c$(name)/:id -  c$(name).show	display a specific item" >> ./controller/c$(name)/c$(name).go
	echo "// GET /c$(name):id/edit -  c$(name).edit	return an HTML form for editing a item" >> ./controller/c$(name)/c$(name).go
	echo "// PUT /c$(name):id -  c$(name).update	update a specific item" >> ./controller/c$(name)/c$(name).go
	echo "// DELETE /c$(name):id -  c$(name).delete	delete a specific item" >> ./controller/c$(name)/c$(name).go

	echo "func Index(rw http.ResponseWriter, r *http.Request, p httprouter.Params){}" >> ./controller/c$(name)/c$(name).go
	echo "func New(rw http.ResponseWriter, r *http.Request, p httprouter.Params){}" >> ./controller/c$(name)/c$(name).go
	echo "func Create(rw http.ResponseWriter, r *http.Request, p httprouter.Params){}" >> ./controller/c$(name)/c$(name).go
	echo "func Show(rw http.ResponseWriter, r *http.Request, p httprouter.Params){}" >> ./controller/c$(name)/c$(name).go
	echo "func Edit(rw http.ResponseWriter, r *http.Request, p httprouter.Params){}" >> ./controller/c$(name)/c$(name).go
	echo "func Upadte(rw http.ResponseWriter, r *http.Request, p httprouter.Params){}" >> ./controller/c$(name)/c$(name).go
	echo "func Delete(rw http.ResponseWriter, r *http.Request, p httprouter.Params){}" >> ./controller/c$(name)/c$(name).go
