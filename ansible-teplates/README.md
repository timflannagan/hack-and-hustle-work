# Purpose

Playing around with jinja2 and Ansible templating adventures.

Goal: Transistioning the metering-operator repository from helm charts to Ansible templating.

## General TODO:

- Determine if static YAML manifests (current stored in the role's `files` directory) can be put in the templates directory
- Determine if using any jinja2 expression needs to have a `.j2` file extension.
- Determine how to organize YAML manifests that only require variable substitution and not jinja2 expressions.
- Determine how to organize any of the `*.tpl` files.

## Charts TODO:

- [x] Metering
- [x] Monitoring
- [ ] Hadoop (mostly done)
- [ ] Presto
- [ ] Hive
- [ ] Ghostunnel
- [ ] Reporting-Operator
- [ ] Openshift-Reporting
