---
draft: "false"
title: "CV"
date: "2022-02-16T18:08:20Z"
resources:
- name: "basic-outline"
- name: "before-tidying"
- name: "makeshift-debugging"
- name: "finished-cv"

---
I recently decided to rewrite my CV, at the start of this process it was still written
in [LaTeX](https://en.wikipedia.org/wiki/LaTeX)
but based on a popular and easy to use LaTeX CV template, [Awesome-CV](https://github.com/posquit0/Awesome-CV), when I
first switched to a LaTeX CV I floundered whilst trying to write my own from scratch and didn't have the inclination to
lean LaTex so went with this popular template.

I've been using this for a number of years, and I've made some fairly substantial tweaks over time to change some
layouts and update fonts to this and deviated from the upstream version considerable, but the more I hack away at it the
harder it was to make Awesome-CV do what I wanted. Everytime I made some minor changes, some very annoying side effects
would crop up, the LaTeX class for this is fairly complicated and quite hard to jump in and edit as you need to, it's
mostly designed to be completely adjusted with variables in the document, so decided I should look at alternatives.

<!--more-->

After looking around at a number of other templates I decided the best way forward, armed with my additional LaTeX
knowledge, would be to create my own from scratch, I had picked up at least the basics of LaTeX from all the tweaks I'd
made and had some time to do some reading on the subject. I started a [new repository](https://github.com/greboid/cv)
and got to work. LaTeX is definitely not user-friendly and documentation for it, ironically, is awful.

I was, however, quite surprised with how quickly I made progress, it only took about an hour to get the basics of a CV
up and running.

{% image './basic-outline.jpg' %}

Obviously this needs a huge amount of work to make it into a CV, and a lot of this is fighting with latex defaults,
whitespace and fonts and margins to get something less ugly.  Unfortunately I didn't make any git commits doing this, so
I can't look back, or demonstrate the process, which is a shame.

Once the basic style was in place the work however wasn't done and there was still a fair amount of work styling it to
look how I feel a CV should look, I took inspiration from Awesome-CV and used this as a base for what I roughly wanted 
it to look like. A fair bit of this was working out how to do multiple columns, this ended up being mini tables all 
over the place, which from my experience with websites feels like its going to come back and bite me someday. I then 
spent a couple of hours trying to debug a minor whitespace issue, latex isn't good at debugging in general and I could
find no way to debug whitespace at all, so I turned to adding horizontal lines around places I was having issues.

{% image './makeshift-debugging.png' %}

One of these was an awkward gap above a list that I'd spent about 2 hours changing various things, and eventually gave
up and pasted the above screenshot and the below source snippet to IRC, asking if anyone could see why there was a
larger gap than I could account for above the blue line (#5) and after the text under the green line (#4).

```latex
\end{tabularx}
# 4 \vspace{0.5em} \\
# 5
\vspace{0.5em}
```

The response was an exceptionally sarcastic "have you tried removing that \\\\", obviously I think I had tried removing
the newline and this not working how I wanted, but I removed the now glaringly obvious new line and my whitespace issue
was solved 🤦.

{% image './finished-cv.jpg' %}

I now have a CV that is nice and simple to update and make changes to, It builds as a PDF, so I can both print this
easily and provide electronically as required, so meets all my initial requirements. As with all my infrastructure these
days I also require version control and automated builds, this is entirely done through gitops and containers, but I'll
save that for another post.
