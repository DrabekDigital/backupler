# Backupler

Backupler is a tool for removing some of the historic backups in given directory according to keep policies
which can allow you to delete old backups within some time window while keeping others in the time window.

This is particulary useful when doing full backups and want to save space on your storage (such as NAS).

## Examples

* You backup daily, but after 3 months you want to keep only one backup per week.
* You backup weekly, but after 2 years you want to remove all backups.
* You backup daily, but after 3 months you want to keep only one backup per week and after 1 year you want to keep only 1 per month.

## Configuration

This is example configuration to show you the possibilities, choose and test your own for your unique purposes.

Beware:
- Policy not spanning the whole interval from `0 days` to `infinity` are considered invalid.
- Only `X days` and `infinity` are supported.
- Only strategies `keep`, `dilute` and `delete` are supported.
- The policies has to be sorted for the sake of clarity!

```yaml
backup:
  naming: "yyyyMMdd_HHmmss" # use exact pattern according to https://docs.oracle.com/javase/6/docs/api/java/text/SimpleDateFormat.html#rfc822timezone
policy:
  - from: "0 day" # inclusive
    to: "30 days" # exclusive
    strategy:
      name: keep
  - from: "30 days"
    to: "90 days"
    strategy:
      name: dilute
      window: "7 days"
  - from: "90 days"
    to: "365 days"
    strategy:
      name: dilute
      window: "30 days"
  - from: "365 days"
    to: "1460 days"
    strategy:
      name: dilute
      window: "365 days"
  - from: "1460 days"
    to: "infinity"
    strategy:
      name: delete
```

## Exucution

```
./backupler \
  --approval
  --test-run \
  --config backups/all/.backupler.yaml \
  backups/all
```

## Limitations

- When diluting is done the **oldest** backup is always kept within the time window!
- **No nesting is supported**, only folders within the provided directory are considered!
- **Only folders matching** the pattern are considered!
- **Only folder matchin dates in the past** are considered!


## Testing

For testing there is simple tool that mocks backup directory for each day in specified interval, run with

```
./backupler \
 --mock-dirs "2018-01-01:2022-12-31:yyyyMMdd_HHmmss" \
 backups/testing
```

## Disclaimer

Always test your configuration in your own playground under real-world scenario. Fiddle around with the exact folder structure (without content) with test mode or in approval mode.

**Author of this library does not provide any guarantee nor is any liable that the tool will work properly and you will do not loose any backup you do not want to loose!**