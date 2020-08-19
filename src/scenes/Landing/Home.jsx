import React from 'react';
import styles from './Home.module.scss';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import checkCircle from '@fortawesome/fontawesome-free-solid/faCheckCircle';
import { Button } from '@trussworks/react-uswds';

const NumberCircle = ({ num }) => <div className={styles['number-circle']}>{num}</div>;

const Step = ({
  actionBtnDisabled,
  actionBtnLabel,
  complete,
  containerClassName,
  description,
  editLabel,
  onEditClick,
  legend,
  onActionBtnClick,
  step,
}) => (
  <div className={containerClassName}>
    <div className={styles['step-header-container']}>
      <h4 className={styles['step-header']}>
        {complete ? (
          <FontAwesomeIcon aria-hidden className={styles.accept} icon={checkCircle} />
        ) : (
          <NumberCircle num={step} />
        )}
        {legend}
      </h4>
      {editLabel && (
        <a href="#" onClick={onEditClick}>
          {editLabel}
        </a>
      )}
    </div>

    <p>{description}</p>
    {actionBtnLabel && (
      <Button disabled={actionBtnDisabled} onClick={onActionBtnClick}>
        {actionBtnLabel}
      </Button>
    )}
  </div>
);

const Home = () => (
  <div style={{ marginTop: -20 }}>
    <header className={styles['customer-header']}>
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
      onEditClick={() => console.log('what')}
    />

    <Step
      containerClassName={styles['step-container']}
      step="2"
      legend="Upload orders"
      description="Upload photos of each page, or upload a PDF"
      actionBtnLabel="Add orders"
      actionBtnDisabled
      onActionBtnClick={() => console.log('some action')}
    />

    <Step
      containerClassName={styles['step-container']}
      step="3"
      legend="Shipments"
      description="Tell us where you're going and when you want to get there. We'll help you set up shipments to make it work"
      actionBtnLabel="Plan your shipments"
      actionBtnDisabled
      onActionBtnClick={() => console.log('some action')}
    />

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
