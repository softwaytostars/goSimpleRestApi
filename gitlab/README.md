With gitlabcCI, it's possible to automatically pull from external github repo to gitlab repo with github actions anf gitlab mirror configuration but one has to pay for this service.

One solution, in order to not spend any money, then could be to push in the same time to github and gitlab.
After creating projet on gitlab from an external github repo:

on local do:

 `git remote add origin https://github.com/softwaytostars/goSimpleRestApi`
 
 `git remote set-url --add --push origin https://gitlab.com/AlexisCot/goSimpleRestApi`

 check all is ok:

 `git remote -v`

 Add your public ssh key to the gitlab project, then:

 on local do:

 `git config --global url.ssh://git@gitlab.com/.insteadOf https://gitlab.com/`

 Then  `git remote -v` gives:

 `origin	https://github.com/softwaytostars/goSimpleRestApi (fetch)
 origin	ssh://git@gitlab.com/AlexisCot/goSimpleRestApi (push)
 origin	https://github.com/softwaytostars/goSimpleRestApi (push)
`

Then

`git push origin branch_name` will push to both repositories

https://gitlab.com/AlexisCot/goSimpleRestApi