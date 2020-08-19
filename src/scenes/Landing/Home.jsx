import React from 'react';
import styles from './Home.module.scss';
import { Button } from '@trussworks/react-uswds';
import { ReactComponent as DocsIcon } from 'shared/icon/documents.svg';
import { ReactComponent as AcceptIcon } from 'shared/icon/accept.svg';
import { ReactComponent as EditIcon } from 'shared/icon/edit.svg';

const NumberCircle = ({ num }) => <div className={styles['number-circle']}>{num}</div>;

const Step = ({
  actionBtnDisabled,
  actionBtnLabel,
  children,
  complete,
  containerClassName,
  description,
  editDisabled,
  editLabel,
  legend,
  onActionBtnClick,
  onEditClick,
  secondary,
  step,
}) => (
  <div className={containerClassName}>
    <div className={styles['step-header-container']}>
      <h4 className={styles['step-header']}>
        {complete ? <AcceptIcon aria-hidden className={styles.accept} /> : <NumberCircle num={step} />}
        {legend}
      </h4>
      {editLabel && (
        <a href={editDisabled ? null : '#'} onClick={editDisabled ? null : onEditClick}>
          {editLabel}
        </a>
      )}
    </div>

    <p>{children ? children : description}</p>
    {actionBtnLabel && (
      <Button disabled={actionBtnDisabled} secondary={secondary} onClick={onActionBtnClick}>
        {actionBtnLabel}
      </Button>
    )}
  </div>
);

const FilesUploaded = ({ files }) => (
  <div className={styles['doc-list-container']}>
    <h6>{files.length} FILES UPLOADED</h6>
    {files.map((file) => (
      <div>
        <DocsIcon className={styles['docs-icon']} />
        {file.filename}
      </div>
    ))}
  </div>
);

const ShipmentListItem = ({ shipment }) => {
  const shipmentClassName = styles[`shipment-list-item-${shipment.shipmentType}`];
  return (
    <div className={`${styles['shipment-list-item-container']} ${shipmentClassName}`}>
      <strong>{shipment.shipmentType}</strong> <span>{shipment.id}</span> <EditIcon />
    </div>
  );
};

const ShipmentList = ({ shipments }) => (
  <div>
    {shipments.map((shipment) => (
      <ShipmentListItem shipment={shipment} />
    ))}
  </div>
);

const shipments = [
  { shipmentType: 'PPM', id: '#123ABC-001' },
  { shipmentType: 'HHG', id: '#123ABC-002' },
  { shipmentType: 'NTS', id: '#123ABC-003' },
];
const Home = () => (
  <div>
    <header className={`${styles['customer-header']} padding-top-3 padding-bottom-3`}>
      <h2>Riley Baker</h2>
      <p>
        You're leaving <strong>Buckley AFB</strong>
      </p>
    </header>
    <Step
      containerClassName={styles['step-container']}
      step="1"
      legend="Profile complete"
      description="Make sure to keep your personal information up to date during your move"
      onActionBtnClick={() => console.log('some action')}
      editLabel="Edit"
      editDisabled
      onEditClick={(e) => {
        e.preventDefault();
        console.log('what');
      }}
    />

    <Step
      containerClassName={styles['step-container']}
      step="2"
      legend="Upload orders"
      description="Upload photos of each page, or upload a PDF"
      actionBtnLabel="Add orders"
      actionBtnDisabled
      onActionBtnClick={() => console.log('some action')}
    >
      <FilesUploaded files={[{ filename: 'file 1' }, { filename: 'file 2' }, { filename: 'file 3' }]} />
    </Step>

    <Step
      containerClassName={styles['step-container']}
      step="3"
      legend="Shipments"
      description="Tell us where you're going and when you want to get there. We'll help you set up shipments to make it work"
      actionBtnLabel="Plan your shipments"
      secondary={shipments.length > 0}
      onActionBtnClick={() => console.log('some action')}
    >
      <ShipmentList shipments={shipments} />
    </Step>

    <Step
      containerClassName={styles['step-container']}
      step="4"
      legend="Confirm move request"
      description="Review your move details and sign the legal paperwork, then send the info on to your move counselor"
      actionBtnLabel="Review and submit"
      actionBtnDisabled
      onActionBtnClick={() => console.log('some action')}
    />

    <div className={styles['footer-container']}>
      <h6>CONTACTS</h6>
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

export default Home;
