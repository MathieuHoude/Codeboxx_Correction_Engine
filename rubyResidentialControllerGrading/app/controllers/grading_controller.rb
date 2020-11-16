class GradingController < ApplicationController
    require 'json'

    def gradeProject
        data = File.read("rspec-results.json")
        json = JSON.parse(data)

        results = {
            :testResults => []
        }

        json["examples"].each do |example| 
            testResult = {
                :description => example["description"],
                :status => example["status"]
            }

            results[:testResults] << testResult
        end

        results

        render json: results.to_json
    end

end
