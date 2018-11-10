<img src="https://diskgem.info/img/diskgem.svg" alt="DiskGem" width="128" />

### What is DiskGem?
DiskGem® is software for secure file transfer over SFTP. 
DiskGem currently offers an easy to use, stable command-line user interface that supports parallel file transfers and other useful features. DiskGem will soon also support creating encrypted archives on the server which offer encryption of stored files as well as metadata obfuscation.

<img src="https://diskgem.info/img/window.png" alt="DiskGem Window" width="400" />

### Using DiskGem
```sh
git clone https://github.com/kaepora/diskgem.git
cd diskgem
make
sudo make install
```

Or, for the truly hurried:

```sh
echo "git clone https://github.com/kaepora/diskgem.git;cd diskgem;make;sudo make install"|sh
```

`diskgem` will launch DiskGem. `man diskgem` will show the manual page.

### Questions and Answers

- **Question:** _Will DiskGem support Microsoft Windows?_
- **Answer:** DiskGem currently is not available natively on Windows due to `cmd.exe` not supporting unicode, 256-color or emojis. However, DiskGem is already fully supported on Windows via the Windows Subsystem for Linux (WSL) which, incidentally, is also used entirely to write DiskGem itself.

- **Question:** _Could you suggest some music suitable for discovering this software?_
- **Answer:** Yes. [Explain](https://oneohtrixpointnever1.bandcamp.com/track/explain) by Oneohtrix Point Never.

- **Question:** _I have opinions about this work that I am fundamentally certain are superior to the views of the author's. I think the author should know my opinions, preferably by my expressing them in a belligerent and demanding tone!_
- **Answer:** Please forget that you ever saw this software and also that we share the same planet.

### Author
Copyright (c) 2018 [Nadim Kobeissi](https://nadim.computer) and released under the MIT License.
