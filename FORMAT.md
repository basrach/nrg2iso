# NRG (file format)

| Filename extension            | .nrg                                                                                                |
| ----------------------------- | --------------------------------------------------------------------------------------------------- |
| Uniform Type Identifier (UTI) | com.nero.nrg-image                                                                                  |
| UTI conformation              | public.iso-image, com.apple.disk-image, public.archive, public.data, public.item, public.disk-image |
| Developed by                  | Nero AG                                                                                             |
| Type of format                | disk image                                                                                          |
| Container for                 | filesystem and volumes                                                                              |

An **NRG file** is a proprietary optical disc image file format originally created by Nero AG for the Nero Burning ROM utility.[1] It is used to store disc images. Other than Nero Burning ROM, however, a variety of software titles can use these image files. For example, Alcohol 120%, or Daemon Tools can mount NRG files onto virtual drives for reading.

Contrary to popular belief, NRG files are not ISO images with a `.nrg` extension and a header attached. They can store audio tracks for Audio CDs, which ISO images cannot. Nero's NRG format is one of the few formats besides BIN/CUE, Alcohol 120%'s MDF/MDS and CloneCD's CCD/IMG/SUB disc image formats to support Mixed Mode CDs which contain audio CD tracks as well as data tracks.

## File format

The file format specification below is unofficial and as such is lacking some data. There may also be errors.

The NRG file format uses a variation of the Interchange File Format (IFF) and stores data in a chain of "chunks". All integer values are stored unsigned in big endian byte order. Version 1 NRG format stores values as 32-bit integers. Nero Burning ROM v5.5 introduced a new NRG file format, version 2, with support for 64-bit integers.

### Header

The NRG format does not store its data as a header at the beginning of a file. It is instead attached at the end of the file like a footer. Image information is stored as a serialized chain of IFF chunks. To get the offset of the first chunk one must read the NRG footer from the last 8 or 12 bytes of the file.

**Nero footer (Version 1)**

| Size (bytes) | Type     | Value or purpose                     |
| ------------ | -------- | ------------------------------------ |
| 4            | Chunk ID | "NERO"                               |
| 4            | 32-bit   | Offset of the first chunk data chain |

**Nero footer (Version 2)**

| Size (bytes) | Type     | Value or purpose                     |
| ------------ | -------- | ------------------------------------ |
| 4            | Chunk ID | "NER5"                               |
| 8            | 64-bit   | Offset of the first chunk data chain |

### Chunks

#### (CUES) Cue Sheet

Available in all versions of NRG file format.

The CUEX chunk is the concatenation of fixed-size blocks, each one representing a cue point.

The index0 points are present even when they are identical to the index1 ones. The index0 points in audio tracks are incorrect if Nero has been asked to record all the sub-channel data (in that case the sector size is 2448 bytes). No index other than 0 or 1 has been encountered, although the chunk format allows for such cue points to be recorded; thus the number of cue blocks seems to be always 2\*(#track + 1): two indices for each track, an index0 for the lead-in and an index1 for the lead-out.

**Version 1**

| Size (bytes) | Type     | Value / Purpose    |
| ------------ | -------- | ------------------ |
| 4            | Chunk ID | "CUES"             |
| 4            | 32-bit   | Chunk size (bytes) |

**Version 2**

| Size (bytes) | Type     | Value / Purpose                                                                            |
| ------------ | -------- | ------------------------------------------------------------------------------------------ |
| 4            | Chunk ID | "CUEX"                                                                                     |
| 4            | 32bit    | Chunk size (bytes)                                                                         |
| 1            | 8bit     | Mode (values found: 0x01 for audio; 0x21 for non copyright-protected audio; 0x41 for data) |
| 1            | 8bit     | Track number (BCD coded; 0xAA for the lead-out area)                                       |
| 1            | 8bit     | Index number (probably BCD coded)                                                          |
| 1            | 8bit     | Padding? (always zero found)                                                               |
| 4            | 32bit    | LBA position in sectors (signed integer value)                                             |

#### (DAOI) DAO Information

Available in all versions of NRG file format.

DAOI chunks store disc at once sessions specific information in two parts. The first part contains data that is specific for the session only. The second part repeats track specific information (grey) once for each track. Parse the SINF chunks to get the number of tracks for a specific session.

**Version 1**

| Size (bytes) | Type     | Value / Purpose                  |
| ------------ | -------- | -------------------------------- |
| 4            | Chunk ID | "DAOI"                           |
| 4            | 32bit    | Chunk size (bytes) big endian    |
| 4            | 32bit    | Chunk size (bytes) little endian |
| 14           | UPC      |                                  |
| 4            | 32bit    | Toc type                         |
| 1            | 8bit     | First track                      |
| 1            | 8bit     | Last track                       |
| 12           | ISRC     |                                  |
| 4            | 32bit    | Sector size                      |
| 4            | 32bit    | Mode                             |
| 4            | 32bit    | Index0 (Pre gap)                 |
| 4            | 32bit    | Index1 (Start of track)          |
| 4            | 32bit    | Index2 (End of track + 1)        |

**Version 2**

| Size (bytes) | Type     | Value / Purpose                                                                                                                                                                                                                                                                                                  |
| ------------ | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 4            | Chunk ID | "DAOX"                                                                                                                                                                                                                                                                                                           |
| 4            | 32bit    | Chunk size (bytes) big endian                                                                                                                                                                                                                                                                                    |
| 4            | 32bit    | Chunk size (bytes) (big endian already encountered; maybe also little endian on some machines)                                                                                                                                                                                                                   |
| 13           | Text     | UPC (or NULLs)                                                                                                                                                                                                                                                                                                   |
| 1            | 8-bit    | Padding? (always NULL found)                                                                                                                                                                                                                                                                                     |
| 2            | 16-bit   | Toc type (values already found: 0x0000 for audio; 0x0001 for data; 0x2001 for Mode 2/form 1 data)                                                                                                                                                                                                                |
| 1            | 8bit     | First track in the session                                                                                                                                                                                                                                                                                       |
| 1            | 8bit     | Last track in the session                                                                                                                                                                                                                                                                                        |
| 12           | Text     | ISRC (or NULLs)                                                                                                                                                                                                                                                                                                  |
| 2            | 16bit    | Sector size in the image file (bytes)                                                                                                                                                                                                                                                                            |
| 2            | 16bit    | Mode of the data in the image file (values already found: 0x0700 for audio; 0x1000 for audio with sub-channel; 0x0000 for data; 0x0500 for raw data; 0x0f00 for raw data with sub-channel; 0x0300 for Mode 2 Form 1 data; 0x0600 for raw Mode 2/form 1 data; 0x1100 for raw Mode 2/form 1 data with sub-channel) |
| 2            | 16bit    | Unknown (always 0x0001 found)                                                                                                                                                                                                                                                                                    |
| 8            | 64bit    | Index0 (Pre-gap) (bytes)                                                                                                                                                                                                                                                                                         |
| 8            | 64bit    | Index1 (Start of track) (bytes)                                                                                                                                                                                                                                                                                  |
| 8            | 64bit    | End of track + 1 (bytes)                                                                                                                                                                                                                                                                                         |

#### (CDTX) CD-text

Available in version 2 NRG file format.

The CDTX chunk is the concatenation of raw CD-text packs of 18 bytes each.

**Version 2**

| Size (bytes) | Type     | Value / Purpose             |
| ------------ | -------- | --------------------------- |
| 4            | Chunk ID | "CDTX"                      |
| 4            | 32bit    | Chunk size (bytes)          |
| 1            | 8bit     | Pack type                   |
| 1            | 8bit     | Pack type (track number)    |
| 1            | 8bit     | Pack number in the block    |
| 1            | 8bit     | Block number etc.           |
| 12           | Text     | NULL-separated text strings |
| 2            | 16bit    | CRC                         |

#### (ETNF) Extended Track Information

Available in all versions of NRG file format.

ETNF chunks are used to store track information for track at once sessions. The data is repeated once for each track. Parse the SINF chunks to get the number of tracks for a specific session.

**Version 1**

| Size (bytes) | Type     | Value / Purpose       |
| ------------ | -------- | --------------------- |
| 4            | Chunk ID | "ETNF"                |
| 4            | 32bit    | Chunk size (bytes)    |
| 4            | 32bit    | Track offset in image |
| 4            | 32bit    | Track length (bytes)  |
| 4            | 32bit    | Mode                  |
| 4            | 32bit    | Start lba on disc     |
| 4            | ?        |                       |

**Version 2**

| Size (bytes) | Type     | Value / Purpose                                                           |
| ------------ | -------- | ------------------------------------------------------------------------- |
| 4            | Chunk ID | "ETN2"                                                                    |
| 4            | 32bit    | Chunk size (bytes)                                                        |
| 8            | 64bit    | Track offset in image                                                     |
| 8            | 64bit    | Track length (bytes)                                                      |
| 4            | 32bit    | Mode (values found: 0x7 for audio; 0x0 for data; 0x3 for Mode 2 data)     |
| 4            | 32bit    | Start lba on disc (sectors) (the start is after a lead-in of 150 sectors) |
| 4            | 32bit    | Unknown (only zero found)                                                 |
| 4            | 32bit    | Track length (sectors)                                                    |

#### (SINF) Session Information

Available in all versions of NRG file format.

Session information chunks should be used to quickly scan the image for session and track count. SINF chunks are always listed in sequential order corresponding to the sessions order. To get more details information about a specific session one must parse the corresponding DAOI or ETNF chunk.

**Version 1 and 2**

| Size (bytes) | Type     | Value / Purpose     |
| ------------ | -------- | ------------------- |
| 4            | Chunk ID | "SINF"              |
| 4            | 32bit    | Chunk size (bytes)  |
| 4            | 32bit    | # tracks in session |

#### (MTYP) Media Type?

Available in all versions of NRG file format.

This chunk and its use is unknown. A value of 1 (big endian) was found in images of several CD (audio or data; CD-ROM or CD-R).

**Version 1 and 2**

| Size (bytes) | Type     | Value / Purpose    |
| ------------ | -------- | ------------------ |
| 4            | Chunk ID | "MTYP"             |
| 4            | 32bit    | Chunk size (bytes) |
| 4            | ?        |                    |

#### (DINF) Disc Information?

Found in TAO images in version 2 of NRG file format. Found in DAO images in version of NRG file format only if Nero was asked not to close the disc.

This chunk and its use is unknown.

**Version 2 (and 1?)**

| Size (bytes) | Type     | Value / Purpose                                         |
| ------------ | -------- | ------------------------------------------------------- |
| 4            | Chunk ID | "DINF"                                                  |
| 4            | 32bit    | Chunk size (bytes)                                      |
| 4            | 32bit    | Unknown (found 0x1 for an unclosed disc; 0x0 otherwise) |

#### (TOCT) TOC T?

Found in TAO images in version 2 of NRG file format.

This chunk and its use is unknown.

**Version 2 (and 1?)**

| Size (bytes) | Type                        | Value / Purpose    |
| ------------ | --------------------------- | ------------------ |
| 4            | Chunk ID                    | "TOCT"             |
| 4            | 32bit                       | Chunk size (bytes) |
| 2            | Unknown (always zero found) |

#### (RELO)

Found in TAO images in version 2 of NRG file format.

This chunk and its use is unknown.

**Version 2 (and 1?)**

| Size (bytes) | Type                        | Value / Purpose    |
| ------------ | --------------------------- | ------------------ |
| 4            | Chunk ID                    | "RELO"             |
| 4            | 32bit                       | Chunk size (bytes) |
| 4            | Unknown (always zero found) |

#### (END!) End of chain

Available in all versions of NRG file format.

End of chain chunk is signals that there are no more chunks to be read.

**Version 1 and 2**

| Size (bytes) | Type     | Value / Purpose          |
| ------------ | -------- | ------------------------ |
| 4            | Chunk ID | "END!"                   |
| 4            | 32bit    | Chunk size (always zero) |
