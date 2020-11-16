FROM ruby:2.5
ARG githubHandle
RUN echo githubHandle
ENV GITHUBHANDLE=${githubHandle}
ARG ssh_prv_key
ENV PRIVATEKEY=${ssh_prv_key}
ARG ssh_pub_key
ENV PUBLICKEY=${ssh_pub_key}

RUN apt-get update -qq && apt-get install -y nodejs yarn
# Authorize SSH Host
RUN mkdir -p /root/.ssh && \
    chmod 0700 /root/.ssh && \
    ssh-keyscan github.com > /root/.ssh/known_hosts
# Add the keys and set permissions
RUN echo "$ssh_prv_key" > /root/.ssh/id_rsa && \
    echo "$ssh_pub_key" > /root/.ssh/id_rsa.pub && \
    chmod 600 /root/.ssh/id_rsa && \
    chmod 600 /root/.ssh/id_rsa.pub

RUN apt-get update -qq && apt-get install -y nodejs 
RUN mkdir /gradingApp
WORKDIR /gradingApp
COPY Gemfile /gradingApp/Gemfile
COPY Gemfile.lock /gradingApp/Gemfile.lock
RUN gem install bundler
RUN bundle install
COPY . /gradingApp
COPY correction-script.sh /usr/bin/
RUN chmod +x /usr/bin/correction-script.sh
# RUN "/usr/bin/correction-script.sh"

# Add a script to be executed every time the container starts.
COPY entrypoint.sh /usr/bin/
RUN chmod +x /usr/bin/entrypoint.sh
ENTRYPOINT ["entrypoint.sh"]
EXPOSE 3000

# Start the main process.
CMD ["rails", "server", "-b", "0.0.0.0"]