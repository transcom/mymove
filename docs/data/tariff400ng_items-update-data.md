# Adding `ShipmentLineItem` records to the 400ng table (`tariff400ng_items`)

## Update the [400ng Common Item Code List](https://docs.google.com/spreadsheets/d/1MSkrhLHH9tHGVGN7ELVLkdpg7XRTb3R3I1xd-ZEkCI4/edit#gid=1382174367)

* Check if the "Service Code (L705)" is present in the [spreadsheet](https://docs.google.com/spreadsheets/d/1MSkrhLHH9tHGVGN7ELVLkdpg7XRTb3R3I1xd-ZEkCI4/edit#gid=1382174367)
* If the Service Code the you want to add the 400ng table isn't in this spreadsheet you will have to add it there.
* If there is a reason that the item is not already present and you are having to add it there yourself, add a Google docs comment to the spreadsheet addressing why the new item is being added.

## Update the table `tariff400ng_items` and the appropriate rate tables

* Gather the information needed, the information should be in the [400ng Common Item Code List](https://docs.google.com/spreadsheets/d/1MSkrhLHH9tHGVGN7ELVLkdpg7XRTb3R3I1xd-ZEkCI4/edit#gid=1382174367) if it's not there you will need to add it there.
* Reach out to the #dp3-ask-the-govies Slack channel and get the information needed, if you do not already have it
* Add the new and/or updated item to the database table (`tariff400ng_items`) and the rate table(s)

### Check that the information isn't already in the table

E.g, Query for the `tariff400ng_items.code`

### Add or update the record

* And the record to the database
* (optional) If you are copying from a very similar record you can use the `INSERT INTO, SELECT` command via command line
  E.g.

  ```sql
    -- Example adding 105C into the table while copying from 105A
    INSERT INTO tariff400ng_items
        (code, discount_type, allowed_location, item, measurement_unit_1, created_at, updated_at)
    SELECT '105C', discount_type, 'DESTINATION', 'Full Unpack', measurement_unit_1, NOW(), NOW()
    FROM tariff400ng_items
    WHERE code = '105A';
   ```

This basically is overkill (in this scenario) since there was only 1 column copied into the new record.

* (optional) Or copy and paste row inside your Postgres editor of choice, changing the necessary rows

### Create the [migration](https://github.com/transcom/mymove/blob/master/docs/database.md#migrations)

## Source code

* For an HHG (Household Goods) move, the function that computes and creates shipment line items is [`ComputeShipment`](https://github.com/transcom/mymove/blob/master/pkg/rateengine/rateengine.go#L153)
* For a PPM (Personally Procured Move) move, the function that computes and creates shipment line items is [`ComputePPM`](https://github.com/transcom/mymove/blob/master/pkg/rateengine/rateengine.go#L73)

## Test

* If you are making any changes to `func CreateBaseShipmentLineItems()` you will have to update the `server_test` `line_items_test.go`