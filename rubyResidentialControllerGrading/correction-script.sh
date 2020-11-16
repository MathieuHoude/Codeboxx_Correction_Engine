#!/bin/bash
rm rspec-results.json
git clone git@github.com:$GITHUBHANDLE/Rocket_Elevators_Controllers.git
mv Rocket_Elevators_Controllers/residential_controller.rb lib/
rspec --format j >> rspec-results.json
rm lib/residential_controller.rb
rm -rf Rocket_Elevators_Controllers