#!/bin/bash
build=../../build.sh 

mkdir -p lib

echo building domainfinder...
$build domainfinder

echo building synonyms...
$build chapter4/thesaurus/synonyms ../../domainfinder/lib/synonyms
cp ../thesaurus/synonyms/test.env ./lib/

echo building available...
$build chapter4/available ../domainfinder/lib/available

echo building sprinkle...
$build chapter4/splinkle ../domainfinder/lib/splinkle

echo building coolify...
$build chapter4/coolify ../domainfinder/lib/coolify

echo building domainify...
$build chapter4/domainify ../domainfinder/lib/domainify


