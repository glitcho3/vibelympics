# sanitize.awk
{
    gsub(/[[:cntrl:]]/, "")
    if (length($0) > 50000) next
}

