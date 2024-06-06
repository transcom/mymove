import { React } from 'react';
import { number, bool } from 'prop-types';
import classnames from 'classnames';

import HeaderSection, { sectionTypes } from './HeaderSection';
import styles from './PPMHeaderSummary.module.scss';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { usePPMCloseoutQuery } from 'hooks/queries';
import { formatCustomerContactFullAddress } from 'utils/formatters';

const GCCAndIncentiveInfo = ({ ppmShipmentInfo }) => {
  const { ppmCloseout, isLoading, isError } = usePPMCloseoutQuery(ppmShipmentInfo.id);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const incentives = {
    isAdvanceRequested: ppmShipmentInfo.hasRequestedAdvance,
    isAdvanceReceived: ppmShipmentInfo.hasReceivedAdvance,
    advanceAmountRequested: ppmShipmentInfo.advanceAmountRequested,
    advanceAmountReceived: ppmShipmentInfo.advanceAmountReceived,
    grossIncentive: ppmCloseout.grossIncentive + ppmCloseout.SITReimbursement,
    gcc: ppmCloseout.gcc,
    remainingIncentive: ppmCloseout.remainingIncentive + ppmCloseout.SITReimbursement,
  };

  const incentiveFactors = {
    haulType: ppmCloseout.haulType,
    haulPrice: ppmCloseout.haulPrice,
    haulFSC: ppmCloseout.haulFSC,
    packPrice: ppmCloseout.packPrice,
    unpackPrice: ppmCloseout.unpackPrice,
    dop: ppmCloseout.dop,
    ddp: ppmCloseout.ddp,
    sitReimbursement: ppmCloseout.SITReimbursement,
  };

  return (
    <>
      <hr />
      <HeaderSection
        sectionInfo={{
          type: sectionTypes.incentives,
          ...incentives,
        }}
        dataTestId="incentives"
      />
      <hr />
      <HeaderSection
        sectionInfo={{ type: sectionTypes.incentiveFactors, ...incentiveFactors }}
        dataTestId="incentiveFactors"
      />
    </>
  );
};
export default function PPMHeaderSummary({ ppmShipmentInfo, ppmNumber, showAllFields }) {
  const shipmentInfo = {
    plannedMoveDate: ppmShipmentInfo.expectedDepartureDate,
    actualMoveDate: ppmShipmentInfo.actualMoveDate,
    pickupAddress: ppmShipmentInfo.pickupAddress
      ? formatCustomerContactFullAddress(ppmShipmentInfo.pickupAddress)
      : '—',
    destinationAddress: ppmShipmentInfo.destinationAddress
      ? formatCustomerContactFullAddress(ppmShipmentInfo.destinationAddress)
      : '—',
    miles: ppmShipmentInfo.miles,
    estimatedWeight: ppmShipmentInfo.estimatedWeight,
    actualWeight: ppmShipmentInfo.actualWeight,
  };

  return (
    <header className={classnames(styles.PPMHeaderSummary)}>
      <div className={styles.header}>
        <h3>PPM {ppmNumber}</h3>
        <section>
          <HeaderSection
            sectionInfo={{
              type: sectionTypes.shipmentInfo,
              advanceAmountReceived: ppmShipmentInfo.advanceAmountReceived,
              ...shipmentInfo,
            }}
            dataTestId="shipmentInfo"
          />
        </section>
        {showAllFields && <GCCAndIncentiveInfo ppmShipmentInfo={ppmShipmentInfo} />}
      </div>
    </header>
  );
}

PPMHeaderSummary.propTypes = {
  ppmNumber: number.isRequired,
  showAllFields: bool.isRequired,
};

PPMHeaderSummary.defaultProps = {};

// TODO: Add shape/propType/defaults for incentives and GCC components here.
