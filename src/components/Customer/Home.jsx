/* eslint-disable react/prop-types */
import React from 'react';
import { string, arrayOf, shape, func } from 'prop-types';
import { Button, Alert } from '@trussworks/react-uswds';

import styles from './Home.module.scss';

import { ReactComponent as DocsIcon } from 'shared/icon/documents.svg';
import { ReactComponent as AcceptIcon } from 'shared/icon/accept.svg';
import { ReactComponent as EditIcon } from 'shared/icon/edit.svg';
import Helper from 'components/Customer/Helper';

const NumberCircle = ({ num }) => <div className={styles['number-circle']}>{num}</div>;

NumberCircle.propTypes = {
  num: string.isRequired,
};

const Step = ({
  actionBtnDisabled,
  actionBtnLabel,
  children,
  complete,
  completedHeaderText,
  containerClassName,
  description,
  editDisabled,
  editLabel,
  headerText,
  onActionBtnClick,
  onEditClick,
  secondary,
  step,
}) => {
  const secondaryClassName = styles['usa-button--secondary'];
  return (
    <div className={`${containerClassName} margin-bottom-6`}>
      <div className={`${styles['step-header-container']} margin-bottom-2`}>
        {complete ? <AcceptIcon aria-hidden className={styles.accept} /> : <NumberCircle num={step} />}
        <strong>{complete ? completedHeaderText : headerText}</strong>
        {editLabel && (
          <Button editDisabled className={styles['edit-button']} onClick={onEditClick}>
            {editLabel}
          </Button>
        )}
      </div>

      {children || <p>{description}</p>}
      {actionBtnLabel && (
        <Button
          className={`margin-top-3 ${secondary ? secondaryClassName : ''}`}
          disabled={actionBtnDisabled}
          onClick={onActionBtnClick}
        >
          {actionBtnLabel}
        </Button>
      )}
    </div>
  );
};

const FilesUploaded = ({ files }) => (
  <div className={`${styles['doc-list-container']} padding-left-2 padding-right-2`}>
    <h6 className="margin-top-2 margin-bottom-2">{files.length} FILES UPLOADED</h6>
    {files.map((file) => (
      <div key={file.filename} className={`margin-bottom-2 ${styles['doc-list-item']}`}>
        <DocsIcon className={styles['docs-icon']} />
        {file.filename}
      </div>
    ))}
  </div>
);

FilesUploaded.propTypes = {
  files: arrayOf(shape({ filename: string.isRequired })).isRequired,
};

const ShipmentListItem = ({ shipment, onShipmentClick }) => {
  function handleEnterOrSpace(event) {
    const key = event.which || event.keyCode; // Use either which or keyCode, depending on browser support
    // enter or space
    if (key === 13 || key === 32) {
      onShipmentClick(shipment);
    }
  }
  const shipmentClassName = styles[`shipment-list-item-${shipment.type}`];
  return (
    <div
      className={`${styles['shipment-list-item-container']} ${shipmentClassName} margin-bottom-1`}
      onClick={() => onShipmentClick(shipment)}
      onKeyDown={(event) => handleEnterOrSpace(event)}
      role="button"
      tabIndex="0"
    >
      <strong>{shipment.type}</strong> <span>{shipment.id}</span> <EditIcon className={styles.edit} />
    </div>
  );
};

ShipmentListItem.propTypes = {
  shipment: shape({ id: string.isRequired, type: string.isRequired }).isRequired,
  onShipmentClick: func.isRequired,
};

const ShipmentList = ({ shipments, onShipmentClick }) => (
  <div>
    {shipments.map((shipment) => (
      <ShipmentListItem key={shipment.id} onShipmentClick={onShipmentClick} shipment={shipment} />
    ))}
  </div>
);

ShipmentList.propTypes = {
  shipments: arrayOf(shape({ id: string.isRequired, type: string.isRequired })).isRequired,
  onShipmentClick: func.isRequired,
};

const shipments = [
  { type: 'PPM', id: '#123ABC-001' },
  { type: 'HHG', id: '#123ABC-002' },
  { type: 'NTS', id: '#123ABC-003' },
];

const Home = () => {
  function handleShipmentClick(shipment) {
    console.log('this is the shipment', shipment);
  }

  return (
    <div>
      <header className={`${styles['customer-header']} padding-top-3 padding-bottom-3 margin-bottom-2`}>
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
        containerClassName={styles['step-container']}
        description="Make sure to keep your personal information up to date during your move"
        editDisabled
        editLabel="Edit"
        onActionBtnClick={() => console.log('some action')}
        step="1"
        onEditClick={(e) => {
          e.preventDefault();
          console.log('what');
        }}
      />

      <Step
        complete
        completedHeaderText="Orders uploaded"
        containerClassName={styles['step-container']}
        description="Upload photos of each page, or upload a PDF"
        editLabel="Edit"
        onEditClick={() => console.log('edit button clicked')}
        headerText="Upload orders"
        onActionBtnClick={() => console.log('some action')}
        step="2"
      >
        <FilesUploaded
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
        containerClassName={styles['step-container']}
        description="Tell us where you're going and when you want to get there. We'll help you set up shipments to make it work"
        headerText="Shipments"
        secondary
        step="3"
      >
        <ShipmentList shipments={shipments} onShipmentClick={handleShipmentClick} />
      </Step>

      <Step
        actionBtnDisabled
        actionBtnLabel="Review and submit"
        containerClassName={styles['step-container']}
        description="Review your move details and sign the legal paperwork, then send the info on to your move counselor"
        headerText="Confirm move request"
        onActionBtnClick={() => console.log('some action')}
        step="4"
      />

      <div className={`${styles['footer-container']} padding-top-2 padding-left-2 padding-right-2 padding-bottom-3`}>
        <h6 className="margin-bottom-1">CONTACTS</h6>
        <p>
          <strong>Seymour Johnson AFB</strong>
          <br />
          Origin Transportation Office
          <br />
          (919) 722-5458
        </p>
      </div>
    </div>
  );
};

export default Home;
