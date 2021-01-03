---
title: godbledger
---
# godbledger

Start Here:
[Quickstart](Quickstart)

If you possibly want to use MySQL rather than Sqlite3 as the database then follow this:
[How to setup MySQL as the backend](How-to-Setup-MySQL-as-the-backend)

Then to set up a fully fledged Systemd service so that your server can run on startup:
[TODO:(sean)](404)

{% assign pages_list = site.pages %}
<ul>
{% for node in pages_list %}
{%   if node.title != null %}
{%     if node.layout == "page" %}
<li>
<a class="sidebar-nav-item{% if page.url == node.url %} active{% endif %}" href="{{ site.baseurl }}{{ node.url  | remove_first: '/' }}">{{ node.title }}</a>
</li>
{%     endif %}
{%   endif %}
{% endfor %}
</ul>