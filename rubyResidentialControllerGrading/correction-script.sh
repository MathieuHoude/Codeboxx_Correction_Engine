#!/bin/bash
git clone git@github.com:$GITHUBHANDLE/Rocket_Elevators_Controllers.git
mv Rocket_Elevators_Controllers/residential_controller.rb lib/
rspec --format j
rm lib/residential_controller.rb
rm -rf Rocket_Elevators_Controllers

# git clone git@github.com:$GITHUBHANDLE/Rocket_Elevators_Controllers.git -q 
# mv Rocket_Elevators_Controllers/residential_controller.rb lib/ >> /dev/null
# rspec --format j
# rm lib/residential_controller.rb > /dev/null
# rm -rf Rocket_Elevators_Controllers > /dev/null