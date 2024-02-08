import React from 'react';
import { number, bool } from 'prop-types';
import classnames from 'classnames';

import HeaderSection, { sectionTypes } from './HeaderSection';
import styles from './PPMHeaderSummary.module.scss';

import { usePPMCloseoutQuery } from 'hooks/queries';

export default function PPMHeaderSummary({ ppmShipmentInfo, ppmNumber, showAllFields }) {
  const ppmCloseout = usePPMCloseoutQuery(ppmNumber);
  const shipmentInfo = {
    plannedMoveDate: ppmShipmentInfo.expectedDepartureDate,
    actualMoveDate: ppmShipmentInfo.actualMoveDate,
    actualPickupPostalCode: ppmShipmentInfo.actualPickupPostalCode,
    actualDestinationPostalCode: ppmShipmentInfo.actualDestinationPostalCode,
    miles: ppmShipmentInfo.miles,
    estimatedWeight: ppmShipmentInfo.estimatedWeight,
    actualWeight: ppmShipmentInfo.actualWeight,
  };
  const incentives = {
    isAdvanceRequested: ppmShipmentInfo.hasRequestedAdvance,
    isAdvanceReceived: ppmShipmentInfo.hasReceivedAdvance,
    advanceAmountRequested: ppmShipmentInfo.advanceAmountRequested,
    advanceAmountReceived: ppmShipmentInfo.advanceAmountReceived,
    grossIncentive: ppmCloseout.grossIncentive,
    gcc: ppmCloseout.gcc,
    remainingIncentive: ppmCloseout.remainingIncentive,
  };
  const gccFactors = {
    haulPrice: ppmCloseout.haulPrice,
    haulFSC: ppmCloseout.haulFSC,
    fullPackUnpackCharge: ppmCloseout.packPrice + ppmShipmentInfo.unpackPrice,
  };
  return (
    <header className={classnames(styles.PPMHeaderSummary)}>
      <div className={styles.header}>
        <h3>PPM {ppmNumber}</h3>
        <section>
          <HeaderSection
            sectionInfo={{
              type: sectionTypes.shipmentInfo,
              ...shipmentInfo,
            }}
          />
        </section>
        <hr />
        {showAllFields && (
          <>
            <HeaderSection
              sectionInfo={{
                type: sectionTypes.incentives,
                ...incentives,
              }}
            />
            <hr />
            <HeaderSection sectionInfo={{ type: sectionTypes.gcc, ...gccFactors }} />
          </>
        )}
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
