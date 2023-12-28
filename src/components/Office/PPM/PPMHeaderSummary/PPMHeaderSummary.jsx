import React from 'react';
import { number } from 'prop-types';
import classnames from 'classnames';

import HeaderSection, { sectionTypes } from './HeaderSection';
import styles from './PPMHeaderSummary.module.scss';

import { PPMShipmentShape } from 'types/shipment';

export default function PPMHeaderSummary({ ppmShipment, ppmNumber, showAllFields }) {
  const { actualPickupPostalCode, actualDestinationPostalCode, actualMoveDate } = ppmShipment || {};

  return (
    <header className={classnames(styles.PPMHeaderSummary)}>
      <div className={styles.header}>
        <h3>PPM {ppmNumber}</h3>
        <section>
          <HeaderSection
            sectionInfo={{
              type: sectionTypes.shipmentInfo,
              actualPickupPostalCode,
              actualMoveDate,
              actualDestinationPostalCode,
            }}
          />
        </section>
        <hr />
        {showAllFields && (
          <>
            <HeaderSection sectionInfo={{ type: sectionTypes.incentives }} />
            <hr />
            <HeaderSection sectionInfo={{ type: sectionTypes.gcc }} />
          </>
        )}
      </div>
    </header>
  );
}

PPMHeaderSummary.propTypes = {
  ppmShipment: PPMShipmentShape,
  ppmNumber: number.isRequired,
};

PPMHeaderSummary.defaultProps = {
  ppmShipment: undefined,
};
