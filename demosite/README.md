Demosite
==========
Demosite is a minimal sample site using DM framework.

How to run demosite
---------------
**Import database**

`mysql -u <username> -p demosite < db/db.sql`

**Configurations**

in configs/site.yaml, change database to correct

**Run**

Under demosite/cmd folder
`go run demo.go dm/demosite`

**Visit**

Visit http://localhost:8092


**Admin(to be changed)**
 - Configuration: change database connection in admin/configs/site.yaml
 - Run: under admin/cmd, run `go run demo.go dm/admin`
 - Visit http://localhost:8089

**Build mode**

How to create a new site based on demosite?
----------------
A simple way to create a new website is copy the demosite to a project and do modifications.

**Import clean database**

**Generate content entities**
After configuring configs/contenttype.json, you need to create/update content entities.

Run below command where dm/demosite is the application package name, dm/dm/codegen can be changed based on your current directory.
`go run dm/dm/codegen/contenttypes/gen.go dm/demosite`

**Configurations**

**Run**

**Build**
