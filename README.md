Libraries for NNTP connections, parsing NZB files & parchive repair. Able to decode yEnc binaries.

NNTP
-------
    conn, err := nntp.Dial(nntpServer, nntpSSL)
    if err != nil {
        /* handle error */
    }

    if err = conn.Authenticate(nntpUser, nntpPass); err != nil {
        /* handle error */
    }

    if article, err := conn.DecodedArticle(articleId); err != nil {
        /* handle error */
    } else {
        /* do something with yEnc decoded article */
    }    

NZB
-------
    if n, err := nzb.ReadNzb(filename); err != nil {
        /* handle error */
    } else {
        // generate a queue of articles to be downloaded, with a 
        // default status
        queue := n.GenerateQueue(defaultStatus)
    }

Par2
-------
    // Verification only dumps to stdout, Stat complete, Repairing in progress
    // will read all parchives for the initial .par2
    stat, err := par2.Stat(par2file)    
    par2.Verify(stat)


The MIT License (MIT)
---------------------
Copyright (c) 2015, Michael D. Tighe

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
