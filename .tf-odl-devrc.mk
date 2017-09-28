# This file is designed to assist you with configuring your environment for
# testing the Active Directory provider, and also serves as a catalog of the environment
# variables that are required to test all of the resources in this provider.
#
# This should be copied to ~/.tf-odl-devrc.mk and edited accordingly.
#
# Note that the use of all of this file is not required - environment variables
# can still be completely set from the command line or your existing
# environment. In this case, don't use this file as it may cause conflicts.
#
# NOTE: Remove the annotations from any variables that have them inline, or
# else make will add the whitespace to the variable as well.
#
# The essentials. None of the tests will run if you don't have these.
export ODL_SERVER_IP       ?= server.ip
export ODL_SERVER_PORT     ?= server.domain
export ODL_SERVER_USER     ?= server.user
export ODL_SERVER_PASSWORD ?= changeme

# vi: filetype=make