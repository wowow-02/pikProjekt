# pikProjekt
Paralelno aproksimativno prepoznavanje uzoraka u tekstu


Iz stvarnog svijeta uzet ćemo problem traženja genskih uzoraka unutar očitanja dobivenih od  uređaja za sekvenciranje genoma. Polazišna točka bit će tekstualna datoteka u tzv. FASTA formatu, koja sadrži  realna izlazna očitanja uređaja za sekvenciranje i to u sljedećem obliku:

 

>MCHU - Calmodulin - Human, rabbit, bovine, rat, and chicken

TACGCATAAGCGCCAAAAGCACGCCGGGCGACCATAATGAGAAGATGCTCACCGCCAGCGTCAACAGTAAGTAGACGATGAGAGCGCCAGGGTAAACAGTTTGCACGGGCCGCGACGAAGGCAGTCGATAATGCCGCCATTGCTGGGATGCTCGCCCCCAGACG...

 

>gi|5524211|gb|AAD44166.1| cytochrome b [Elephas maximus maximus]

CGCGTCGTTAATGCCGTCGCCCACCATTGCCACCTGACGTCCTTCACTTTGCAGATGTTATCCAACCTTCATAGTTACTCGTCCGGCTACGCGGTCGCGATCACC...

…

 

Svako očitanje započinje znakom “>” (i kratkim opisom o p) te sadrži niz znakova A, T, C i G te  je potpuno nezavisno o drugima (redoslijed očitanja ne mora odgovarati realnom redoslijedu u stvarnom genomu niti se očitanja moraju nadovezivati jedno na drugo). Duljine očitanja također nisu specificirane i variraju. Zbog te nezavisnosti moguće je paralelno izdvojiti nizove A,T,C,G i ubaciti ih u polje na način da jedan niz/očitanje bude jedan element polja i to bez straha da nam se primjerice dio traženog genskog uzorka nalazi na kraju jednog očitanja, a dio na početku nekog drugog. Sljedeći korak bila bi implementacija algoritama za aproksimativno prepoznavanje uzoraka:

    1.  Algoritam Lavenshtein

    2. Algoritam Rabin-Karp

    3. Naivni algoritam

    4. Algoritam Z

    5. Algoritam Aho-Corasick

    6. Algoritam Baeza-Yates-Gonnet

    7. Algoritam Wu-Manber

    8. Algoritam Bitap (ili Baeza-Yates-Gonnet)

(Odabrali bi 4-5 najzanimljivijih). Svaki algoritam bi izvodili prosljeđujući mu fiksni traženi uzorak znakova A,T,C,G i parametar tolerancije (broj dopuštenih pogrešaka u nađenom podnizu) koji bi varirali. Uzorak bi se tražio (ponovo paralelno) u prije spomenutom polju očitanja i to tako da svakoj dretvi proslijedimo skup od n uzastopnih elemenata polja koja sadrže očitanja uređaja za sekvenciranje. Elemente tog polja prethodno bi rasporedili da osiguramo bolji “load balancing” (jer su očitanja različitih duljina) tako da bi svaku dretvu dopao otprilike jednako dug skup očitanja da na njemu nađe traženi uzorak upotrebom pojedinog algoritma. U osnovi, veći broj dopuštenih pogrešaka trebao bi povećati vrijeme potrebno za pretragu jer algoritam mora provjeriti više mogućih podudaranja. U skladu s time parametar tolerancije trebao bi imati veći utjecaj na vrijeme izvođenja čak i od same duljine traženog uzorka. Međutim, točan utjecaj može se razlikovati ovisno o specifičnom algoritmu i načinu na koji se pogreške obrađuju. Za svaki algoritam bi nacrtali graf ovisnosti vremena izvođenja o iznosu parametra tolerancije što bi nam u konačnici trebalo dati odgovor na pitanje koji je korištenih algoritama je najrobusniji na broj dopuštenih pogrešaka (gledali bi nagibe krivulja).
