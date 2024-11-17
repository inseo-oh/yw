# YJK's silly web adventures (YW)

I've first started learning the Web from the high-level - Writing HTML, CSS, JS, and using various Web APIs.

But how about starting from the... other end?

## Current progress

Note that `main.cc`, containing the `main` function, is mostly playground for trying out APIs I implemented. In other words, it will have whatever random s-it I had last time.

### [DOM](https://dom.spec.whatwg.org/)

For a while this would be my main focus, since parsing HTML requires working DOM.

```c++
std::shared_ptr<dom::Document> document = dom::Document::_create("Document");
std::shared_ptr<dom::Node> html = document->create_element("html");
html->append(document);
```

Currently creating and inserting nodes, and creating elements are supported. Some internal functions for removing nodes are also there though.

## Would this evolve to a full browser?

Maybe. But even then, I don't intend to replicate what Chrome or Firefox does. In other words, I won't be making a browser for regular end-user.


