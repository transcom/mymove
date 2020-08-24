/* eslint-disable no-console */
import React from 'react';
import { Alert } from '@trussworks/react-uswds';

import styles from './Home.module.scss';

import Helper from 'components/Customer/Home/Helper';
import Step from 'components/Customer/Home/Step';
import DocsUploaded from 'components/Customer/Home/DocsUploaded';
import ShipmentList from 'components/Customer/Home/ShipmentList';
import Contact from 'components/Customer/Home/Contact';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const shipments = [
  { type: SHIPMENT_OPTIONS.PPM, id: '#123ABC-001' },
  { type: SHIPMENT_OPTIONS.HHG, id: '#123ABC-002' },
  { type: SHIPMENT_OPTIONS.NTS, id: '#123ABC-003' },
];

const Home = () => {
  function handleShipmentClick(shipment) {
    console.log('this is the shipment', shipment);
  }

  return (
    <div className={`usa-prose grid-container ${styles['grid-container']}`}>
      <header className={styles['customer-header']}>
        <h2>Riley Baker</h2>
        <p>
          You&apos;re leaving <strong>Buckley AFB</strong>
        </p>
      </header>
      <Alert className="margin-top-2 margin-bottom-2" slim type="success">
        Thank you for adding your Profile information
      </Alert>

      <Helper
        title="Next step: Add your orders"
        helpList={[
          'If you have a hard copy, you can take photos of each page',
          'If you have a PDF, you can upload that',
        ]}
      />
      <Step
        complete
        completedHeaderText="Profile complete"
        editBtnDisabled
        editBtnLabel="Edit"
        headerText="Profile complete"
        onActionBtnClick={() => console.log('some action')}
        step="1"
        onEditClick={(e) => {
          e.preventDefault();
          console.log('edit clicked');
        }}
      >
        <p className={styles.description}>Make sure to keep your personal information up to date during your move</p>
      </Step>

      <Step
        complete
        completedHeaderText="Orders uploaded"
        editBtnLabel="Edit"
        onEditClick={() => console.log('edit button clicked')}
        headerText="Upload orders"
        onActionBtnClick={() => console.log('some action')}
        step="2"
      >
        <DocsUploaded
          files={[
            { filename: 'Screen Shot 2020-09-11 at 12.56.58 PM.png' },
            { filename: 'Screen Shot 2020-09-11 at 12.58.12 PM.png' },
            { filename: 'orderspage3_20200723.png' },
          ]}
        />
      </Step>

      <Step
        actionBtnLabel="Add another shipment"
        complete
        completedHeaderText="Shipments"
        headerText="Shipments"
        secondaryBtn
        secondaryClassName="margin-top-2"
        step="3"
      >
        <ShipmentList shipments={shipments} onShipmentClick={handleShipmentClick} />
      </Step>

      <Step
        actionBtnDisabled
        actionBtnLabel="Review and submit"
        containerClassName="margin-bottom-8"
        headerText="Confirm move request"
        onActionBtnClick={() => console.log('some action')}
        step="4"
      >
        <p className={styles.description}>
          Review your move details and sign the legal paperwork, then send the info on to your move counselor
        </p>
      </Step>
      <Contact
        header="Contacts"
        dutyStationName="Seymour Johnson AFB"
        officeType="Origin Transportation Office"
        telephone="(919) 722-5458"
      />
    </div>
  );
};

export default Home;
