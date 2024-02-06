import React from 'react';
import { number, bool } from 'prop-types';
import classnames from 'classnames';

import HeaderSection, { sectionTypes } from './HeaderSection';
import styles from './PPMHeaderSummary.module.scss';

import { PPMCloseoutShape } from 'types/shipment';

export default function PPMHeaderSummary({ ppmCloseout, ppmNumber, showAllFields }) {
  const shipmentInfo = {
    plannedMoveDate: ppmCloseout.plannedMoveDate,
    actualMoveDate: ppmCloseout.actualMoveDate,
    actualPickupPostalCode: ppmCloseout.actualPickupPostalCode,
    actualDestinationPostalCode: ppmCloseout.actualDestinationPostalCode,
    miles: ppmCloseout.miles,
    estimatedWeight: ppmCloseout.estimatedWeight,
    actualWeight: ppmCloseout.actualWeight,
    advanceRequested: ppmCloseout.advanceRequested,
    advanceReceived: ppmCloseout.advanceReceived,
    aoa: ppmCloseout.aoa,
  };
  const incentives = {
    grossIncentive: ppmCloseout.grossIncentive,
    gcc: ppmCloseout.gcc,
    remainingIncentive: ppmCloseout.remainingIncentive,
  };
  const gccFactors = {
    lineHaulPrice: ppmCloseout.lineHaulPrice,
    lineHaulFuelSurcharge: ppmCloseout.lineHaulFuelSurcharge,
    shorthaulPrice: ppmCloseout.shorthaulPrice,
    shorthaulFuelSurcharge: ppmCloseout.shorthaulFuelSurcharge,
    fullPackUnpackCharge: ppmCloseout.packPrice + ppmCloseout.unpackPrice,
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
  ppmCloseout: PPMCloseoutShape,
  ppmNumber: number.isRequired,
  showAllFields: bool.isRequired,
};

PPMHeaderSummary.defaultProps = {
  ppmCloseout: undefined,
};

// TODO: Add shape/propType/defaults for incentives and GCC components here.
