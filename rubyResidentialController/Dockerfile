FROM ruby:2.5
ARG githubHandle
RUN apt-get update -qq && apt-get install -y nodejs yarn
RUN mkdir /app
WORKDIR /app
COPY Gemfile /app/Gemfile
COPY Gemfile.lock /app/Gemfile.lock
RUN bundle install
COPY . /app
COPY correction-script.sh /usr/bin/
RUN chmod +x /usr/bin/correction-script.sh
ENTRYPOINT ["correction-script.sh", $githubHandle]
