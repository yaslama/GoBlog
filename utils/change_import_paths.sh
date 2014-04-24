
find ../ -iname "*.go" | xargs perl -p -i -E 's"github.com/fuxiaohei/GoBlog"github.com/oyvindsk/GoBlog"g'

