#!/bin/bash
GITHUBHANDLE=$1

git clone git@github.com:$GITHUBHANDLE/Rocket_Elevators_Controllers.git
mv Rocket_Elevators_Controllers/residential_controller.rb lib/
rspec --format j << rspec-results.json
