Editorial User Interface(EUI)
===============

Idea of frontend editing.
--------------
Frontend editing will be similar to toolbar above the list of article/folder. It's not only used for article full, but also for blocks, lines if needed also.

Structure wise, all the toolbars will be hooked into content structure shown in frontend, either in the menu, or in the block, or in the list, or in the top/bottom.

**Senarios:**
- add news in blocks
- change order in blocks
- [need more thinking]move by drag-drop
- add more picture in slideshow
- multi upload
- add more in footer, header contents
Most of the implementation will be able to reuse implementation in backend.


**The principle of eui frontend**
Everything is based on context. Examples of context:
- content full view
- children/list
- content block
- slideshow(special list/children)

**Extendable**
The ui should be extendable for more features for frontend editing.



**Things that can be done in frontend**
- Content create, edit, delete
- Insert image, relation in content, meaning browsing is possible. Upload image when browsing
- Content order change
- Add/edit block eg. rss, video
- Special content type creating/editing.
    - slideshow/gallery: add new image, drag&drop image directly
    - frontpage for subsite, permission controlled to attribute level if needed.
    - custom modules(eg. contrast), block system
    -

**Current design**
- When you click node in edit mode, you will browse content and can drag and drop to.
