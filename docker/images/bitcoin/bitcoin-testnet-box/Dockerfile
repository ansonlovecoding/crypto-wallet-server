# bitcoin-testnet-box docker image

FROM ubuntu
LABEL maintainer="Sean Lavine <lavis88@gmail.com>"

# install make
RUN apt-get update && \
	apt-get install --yes make wget bash curl

# create a non-root user
RUN adduser --disabled-login --gecos "" tester

# run following commands from user's home directory
WORKDIR /home/tester

ENV BITCOIN_CORE_VERSION "24.0.1"

# download and install bitcoin binaries
RUN mkdir tmp \
	&& cd tmp \
	&& wget "https://bitcoincore.org/bin/bitcoin-core-${BITCOIN_CORE_VERSION}/bitcoin-${BITCOIN_CORE_VERSION}-x86_64-linux-gnu.tar.gz" \
	&& tar xzf "bitcoin-${BITCOIN_CORE_VERSION}-x86_64-linux-gnu.tar.gz" \
	&& cd "bitcoin-${BITCOIN_CORE_VERSION}/bin" \
	&& install --mode 755 --target-directory /usr/local/bin *

# clean up
RUN rm -r tmp

# copy the testnet-box files into the image
ADD . /home/tester

# make tester user own the bitcoin-testnet-box
RUN chown -R tester:tester /home/tester/bitcoin-testnet-box

# color PS1
RUN mv /home/tester/bitcoin-testnet-box/.bashrc /home/tester/ && \
	cat /home/tester/.bashrc >> /etc/bash.bashrc

# use the tester user when running the image
USER tester

# expose two rpc ports for the nodes to allow outside container access
EXPOSE 19001 19011

ENTRYPOINT /home/tester/script/start.sh

# run commands from inside the testnet-box directory
WORKDIR /home/tester/bitcoin-testnet-box

CMD ["/bin/bash"]

