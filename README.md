# YJK's silly web adventures (YW)

I've first started learning the Web from the high-level - Writing HTML, CSS, JS, and using various Web APIs.

But how about starting from the... other end?

## Current progress

Everything is in WIP and untested basically. For now, this repo is mostly for backing up my source code.

## Copyright

- This project is licensed under BSD 3-Clause license. (See `LICENSE`)

- This project includes material from [WHATWG standards](https://whatwg.org/ipr-policy#711-living-standards-and-review-drafts), such as HTML, DOM, Encoding and Infra, for documentation purposes. Copyright © WHATWG (Apple, Google, Mozilla, Microsoft).

> WHATWG Living Standards and Review Drafts are licensed under Creative Commons "Attribution 4.0 International (CC BY 4.0)". To the extent portions of such Living Standards or Review Drafts are incorporated into source code, such portions in the source code are licensed under the BSD 3-Clause License instead. instead.

## Sources of JSON files

- `yw/res/htmlNamedCharRefs.json` came from [HTML standard's "13.5 Named character references"](https://html.spec.whatwg.org/multipage/named-characters.html#named-character-references) ([JSON file](https://html.spec.whatwg.org/entities.json)). However, the standard states that the list will never change, so in the future I might end up just including it as a Lua file instead to avoid parsing it at runtime.
