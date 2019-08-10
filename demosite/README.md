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
`go run demo.go ..`

**Visit**

Visit http://localhost:8092


**Admin(to be changed)**
 - Configuration: change database connection in admin/configs/site.yaml
 - Run: under admin/cmd, run `go run demo.go ..`
 - Visit http://localhost:8089

**Build mode**

How to create a new site based on demosite?
----------------
A simple way to create a new website is copy the demosite to a project and do modifications.

**Import clean database**

**Configurations**

**Run**

**Build**
