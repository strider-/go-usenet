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

NZB support complete
-------
    if n, err := nzb.ReadNzb(filename); err != nil {
        /* handle error */
    } else {
        // generate a queue of articles to be downloaded, with a 
        // default status
        queue := n.GenerateQueue(defaultStatus)
    }

Par2 
support pending: Verification / Reporting complete, Repairing in progress
-------
    // will read all parchives for the initial .par2
    stat, err := par2.Stat(par2file)    
    par2.Verify(stat)