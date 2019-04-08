import React from 'react';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faBan from '@fortawesome/fontawesome-free-solid/faBan';

function AllowableExpenses(props) {
  function goBack() {
    props.history.goBack();
  }

  return (
    <div className="usa-grid allowable-expenses-container">
      <div>
        <a onClick={goBack}>{'<'} Back</a>
      </div>
      <h3 className="title">Allowable expenses:</h3>
      <p>
        <strong>Storage Expenses</strong>
        <br />
        Storage expenses are <strong>reimburseable</strong> for up to 90 days
      </p>

      <p>
        <strong>Moving Expenses</strong>
        <br />
        Claimable moving expenses will <strong>reduce taxes</strong> on your payment.<br />
        <FontAwesomeIcon aria-hidden className="icon" icon={faCheck} />Claimable expenses include:
        <ul>
          <li>Consumable packling materials</li>
          <li>Contracted expenses</li>
          <li>Gas</li>
          <li>Oil</li>
          <li>Rental equipment</li>
          <li>Tolls</li>
          <li>Weighing fees</li>
          <li>Other</li>
        </ul>
        <br />
        <FontAwesomeIcon aria-hidden className="icon" icon={faBan} />Non-claimable expenses include:
        <ul>
          <li>Animal costs (kennels, transportation)</li>
          <li>Ectra drivers</li>
          <li>Gas</li>
          <li>Hitch fees and tow bars</li>
          <li>Locks</li>
          <li>Meals and lodging</li>
          <li>Moving insurance</li>
          <li>Oil change and routine maintenance</li>
          <li>Purchased auto transporters and dollies</li>
          <li>Sales tax</li>
          <li>Tire chains</li>
        </ul>
      </p>
      <p>
        <strong>Gas and fuel expenses</strong>
        <br />
        Fuel expenses may no longer be claimed in conjunction with a PPM unless the amount of fuel exceeds the amount
        paid for Mileage and Per Diem fees on your travel pay. The IRS will not allow you to claim an expense if you
        were paid for the mileage already. Doing so could result in an IRS audit. However, if your fuel costs exceed the
        Per Diem payment you received, you may claim the portion that exceeds that amount.
      </p>
      <p>
        <strong>When are receipts required?</strong>
        <br />
        Receipts are required for any contracted expenses, storage facilities, and for any expense over $75. If you have
        an expense under $75, you will only need a receipt if you have multiple expenses within that same expense
      </p>
      {/* TODO: Copy WIP still*/}
      <p>If you are missing a required receipt, you can fill out a missing/lost receipt form to submit to Finance.</p>
      <div className="usa-grid button-bar">
        <button className="usa-button-secondary" onClick={goBack}>
          Back
        </button>
      </div>
    </div>
  );
}

export default AllowableExpenses;
