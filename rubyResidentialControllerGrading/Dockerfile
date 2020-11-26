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



