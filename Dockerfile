# Use a TeX Live base image
FROM texlive/texlive

# Set up a working directory
WORKDIR /data

RUN tlmgr option repository https://mirror.ctan.org/systems/texlive/tlnet
RUN tlmgr install fancyhdr

# Install XeLaTeX and additional LaTeX packages
RUN apt-get update && apt-get install -y \
    texlive-xetex \
    texlive-latex-base \
    texlive-latex-extra \
    texlive-fonts-recommended

COPY ./fonts /usr/local/share/fonts/custom
RUN fc-cache -f -v

# Command to run pdflatex by default
ENTRYPOINT ["pdflatex"]