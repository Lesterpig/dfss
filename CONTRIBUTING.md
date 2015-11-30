# Contributing to this project

Please take a moment to review this document in order to make the contribution
process easy and effective for everyone involved.

Following these guidelines helps to communicate that you respect the time of
the developers managing and developing this open source project. In return,
they should reciprocate that respect in addressing your issue or assessing
patches and features.

Please use english for public communication.

## Using the issue tracker

The [issue tracker](https://gitlab.insa-rennes.fr/mpcs/dfss/issues) is the preferred channel for [bug reports](#bug-reports),
[features requests](#feature-requests) and [submitting pull requests](#pull-requests),
but please respect the following restrictions:

* Please **do not** use the issue tracker for personal support requests (use
  [Stack Overflow](http://stackoverflow.com) or IRC).

* Please **do not** derail or troll issues. Keep the discussion on topic and
  respect the opinions of others.

## Bug reports

A bug is a _demonstrable problem_ that is caused by the code in the repository.
Good bug reports are extremely helpful - thank you!

Guidelines for bug reports:

1. **Use the issue search** &mdash; check if the issue has already been
   reported.

2. **Check if the issue has been fixed** &mdash; try to reproduce it using the
   latest `master` or development branch in the repository.

A good bug report shouldn't leave others needing to chase you up for more
information. Please try to be as detailed as possible in your report. What is
your environment? What steps will reproduce the issue? What browser(s) and OS
experience the problem? What would you expect to be the outcome? All these
details will help people to fix any potential bugs.

Example:

> Short and descriptive example bug report title
>
> A summary of the issue and the browser/OS environment in which it occurs. If
> suitable, include the steps required to reproduce the bug.
>
> 1. This is the first step
> 2. This is the second step
> 3. Further steps, etc.
>
> Any other information you want to share that is relevant to the issue being
> reported. This might include the lines of code that you have identified as
> causing the bug, and potential solutions (and your opinions on their
> merits).


## Feature requests

Feature requests are welcome. But take a moment to find out whether your idea
fits with the scope and aims of the project. It's up to *you* to make a strong
case to convince the project's developers of the merits of this feature. Please
provide as much detail and context as possible.

Features and bugfixes are referenced and attributed in a
[Taiga Project](https://tree.taiga.io/project/lesterpig-mpcs/kanban).

## Git branches policy

The `master` branch is intended to provide a *"semi-stable"* version of the dfss
application. It means that the code shall be fully reviewed under `master`, yet
providing the last shiny features and bugfixes.

Features and bugfixes shall be developed under their own branches, and then
merged into `master` via a *Pull request*. **Please remove feature branches after
merge.**

As a best practice, every contributor should use the following naming convention
for feature branches:

```
Fixes    : tg_fix_<bug_description>
Features : tg_<feature>

The "tg_" prefix is referred to the corresponding Taiga.io user-story / issue / task,
and can be ommited.

Examples :
95_fix_crypto_tests
158_platform_https_auth
ssl_hotfix
```

You can rebase. But, please, do it locally (when refactoring your own topic branch).
You don't want to have contributors against you after a missed forced push.

## Pull requests

Good pull requests (patches, improvements, new features) are a fantastic
help. They should remain focused in scope and avoid containing unrelated
commits.

**Please ask first** before embarking on any significant pull request (e.g.
implementing features, refactoring code, porting to a different language),
otherwise you risk spending a lot of time working on something that the
project's developers might not want to merge into the project.

Please adhere to the coding conventions used throughout a project (indentation,
accurate comments, etc.) and any other requirements (tests and their coverage).

We are using **Golang** as the main language. Coding conventions are really simple
and can be found [here](https://golang.org/doc/effective_go.html).

Follow this process if you'd like your work considered for inclusion in the
project:

1. Fork the project, clone your fork,
   and configure the remotes:

   ```bash
   # Clone your fork of the repo into the current directory
   git clone https://gitlab.insa-rennes.fr/<your-username>/<repo-name>
   # Navigate to the newly cloned directory
   cd <repo-name>
   # Assign the original repo to a remote called "upstream"
   git remote add upstream https://gitlab.insa-rennes.fr/mpcs/dfss
   ```

2. If you cloned a while ago, get the latest changes from upstream:

   ```bash
   git checkout <dev-branch>
   git pull upstream <dev-branch>
   ```

3. Create a new topic branch (off the main project development branch) to
   contain your feature, change, or fix:

   ```bash
   git checkout -b <topic-branch-name>
   ```

4. Commit your changes in logical chunks. Please adhere to these [git commit
   message guidelines](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)
   or your code is unlikely be merged into the main project. Use Git's
   [interactive rebase](https://help.github.com/articles/interactive-rebase)
   feature to tidy up your commits before making them public.

   You can add a `[prefix]` at the beginning of the commit message if relevant.

5. Locally merge (or rebase) the upstream development branch into your topic branch:

   ```bash
   git pull [--rebase] upstream <dev-branch>
   ```

6. Push your topic branch up to your fork:

   ```bash
   git push origin <topic-branch-name>
   ```

7. Open a Pull Request with a clear title and description.

**IMPORTANT**: By submitting a patch, you agree to allow the project owner to
license your work under the same license as that used by the project.

# Conduct

We are committed to providing a friendly, safe and welcoming environment for
all, regardless of gender, sexual orientation, disability, ethnicity, religion,
or similar personal characteristic.

On IRC, please avoid using overtly sexual nicknames or other nicknames that
might detract from a friendly, safe and welcoming environment for all.

Please be kind and courteous. There's no need to be mean or rude.
Respect that people have differences of opinion and that every design or
implementation choice carries a trade-off and numerous costs. There is seldom
a right answer, merely an optimal answer given a set of values and
circumstances.

Please keep unstructured critics to a minimum. If you have solid ideas you
want to experiment with, make a fork and see how it works.

We will exclude you from interaction if you insult, demean or harass anyone.
That is not welcome behaviour. We interpret the term "harassment" as
including the definition in the
[Citizen Code of Conduct](http://citizencodeofconduct.org/);
if you have any lack of clarity about what might be included in that concept,
please read their definition. In particular, we don't tolerate behavior that
excludes people in socially marginalized groups.

Likewise any spamming, trolling, flaming, baiting or other attention-stealing
behaviour is not welcome.

