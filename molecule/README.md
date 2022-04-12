# Molecule

Instructions from https://molecule.readthedocs.io/en/latest/getting-started.html

# Installation and setup
```sh
# Make sure python, ansible and docker are installed
virtualenv .venv
source .venv/bin/activate
# Best to reinstall ansible, using a combination of molecule installed in virtualenv and system-wide ansible install
# gave me some issues.
pip install ansible molecule 'molecule[docker]' pytest-testinfra==6.7.0

# When creating a new role (it's also possible to onboard an existing role)
# This will generate the role and molecule directory structure
molecule init role acme.my_new_role --driver-name docker
```

# Basic Molecule Usage

```sh
# Create your molecule target test instance (by default this is docker (=driver-name set before))
molecule create
# You can list your instance:
molecule list
docker ps

# You can login to this instance
# NOTE: at this time, there was a bug in molecule that prevented this from working: https://github.com/ansible-community/molecule/pull/3468
# While the bug has been fixed, it's currently unreleased. Easily fixed by downgrading molecule:
# pip install molecule==3.5.2
molecule login

# Apply your role to the target provider (docker by default)
# What's really happning is that molecule/default/converge.yml is executed, which by default just applies the role
molecule converge
# By running `molecule login` again, you can manually verify that everything was applied

# Run your molecule tests against the instance (see molecule/default/verify.yml)
molecule verify

# When unspecified, molecule runs the 'default' scenario. But you can also pass specific scenarios.
# In this example, we use a custom scenario to run testinfra tests
molecule verify -s myscenario

# Cleanup docker instance
molecule destroy

# There's more intermediate steps, like prepare, dependency, side-effect, etc, that you can also run independently
# using `molecule prepare`, `molecule dependency`, etc. Those are useful if you want to further customize your test
# setup.
```

# Running tests sequences

Molecule has a number of scenario runners built-in that will go through the entire cycle of setting up:
- **check_sequence**: dependency->cleanup->destroy->create->prepare->converge->check->destroy
- **test_sequence**:  dependency->lint->cleanup->destroy->syntax->create->prepare->converge->idempotence->side_effect->verify->cleanup->destroy->

This makes it easy to run e2-e test cases.
Note: you can actually configure what this sequence is in your `molecule.yml` file.

```sh
# Run the test sequence. 
molecule test

# Run the check sequence. Note that this runs ansible check mode: `ansible-playbook --check`, after doing a molecule converge
molecule check
```

# Notes
- Linting used to be part of molecule (`molecule lint` but has been removed by default from molecule v3 onwards as
 "The issue was that molecule runs on scenarios and linting is usually performed at repository level."
- The `tests/` directory directly under the role has nothing to do with molecule, see here: 
  https://old.reddit.com/r/ansible/comments/drypmf/what_is_a_roles_test_folder_for/
    -  Note that `molecule init` generates this directory as it just follows the default template for ansible galaxy
       compatible roles.
- Similarly, the `ansible-test` command is something entirely still which is used as part of ansible module development
  i.e. this is not somethign that end-users would use: https://docs.ansible.com/ansible/latest/dev_guide/testing_units.html