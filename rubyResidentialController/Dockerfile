FROM ruby:2.5
ARG githubHandle
RUN apt-get update -qq && apt-get install -y nodejs yarn
RUN git clone rubyGrading
WORKDIR /rubyGrading
RUN bundle install
COPY . /rubyGrading
COPY correction-script.sh /usr/bin/
RUN chmod +x /usr/bin/correction-script.sh
ENTRYPOINT ["correction-script.sh", $githubHandle]
