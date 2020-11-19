# FROM ruby:2.5

# RUN apt-get update 
# # Authorize SSH Host
# RUN mkdir -p /root/.ssh && \
#     chmod 0700 /root/.ssh && \
#     ssh-keyscan github.com > /root/.ssh/known_hosts
# RUN echo "LogLevel ERROR" > /etc/ssh/ssh_config

# RUN mkdir /gradingApp
# WORKDIR /gradingApp
# COPY Gemfile /gradingApp/Gemfile
# COPY Gemfile.lock /gradingApp/Gemfile.lock
# RUN gem install bundler
# RUN bundle install
# COPY . /gradingApp
# COPY correction-script.sh /usr/bin/
# RUN chmod +x /usr/bin/correction-script.sh

# RUN mkdir lib
# RUN mkdir spec
# RUN mv spec_helper.rb spec
# RUN mv residential_controller_spec.rb spec

# # Start the main process.
# CMD /usr/bin/correction-script.sh


FROM mathieuhoude/ruby-residential-controller-grading:latest
ARG ssh_prv_key
ARG ssh_pub_key
# Add the keys and set permissions
RUN echo "$ssh_prv_key" > /root/.ssh/id_rsa && \
    echo "$ssh_pub_key" > /root/.ssh/id_rsa.pub && \
    chmod 600 /root/.ssh/id_rsa && \
    chmod 600 /root/.ssh/id_rsa.pub

# Start the main process.
CMD /usr/bin/correction-script.sh