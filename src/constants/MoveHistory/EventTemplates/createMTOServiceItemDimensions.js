import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { convertFromThousandthInchToInch } from 'utils/formatters';

export default {
  action: a.INSERT,
  eventName: o.createMTOServiceItem,
  tableName: t.mto_service_item_dimensions,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Requested service item',
  getDetailsLabeledDetails: ({ changedValues }) => {
    const {
      type,
      height_thousandth_inches: heightThousandthInches,
      length_thousandth_inches: lengthThousandthInches,
      width_thousandth_inches: widthThousandthInches,
    } = changedValues;
    const height = convertFromThousandthInchToInch(heightThousandthInches);
    const length = convertFromThousandthInchToInch(lengthThousandthInches);
    const width = convertFromThousandthInchToInch(widthThousandthInches);

    const name = type === 'ITEM' ? 'item_size' : 'crate_size';

    const newChangedValues = {
      ...changedValues,
    };
    newChangedValues[name] = `${height}x${length}x${width} in`;

    return newChangedValues;
  },
};
