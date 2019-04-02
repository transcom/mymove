import React from 'react';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faBan from '@fortawesome/fontawesome-free-solid/faBan';

function AllowableExpenses(props) {
  function goBack() {
    props.history.goBack();
  }

  return (
    <div className="usa-grid">
      <div>
        <a onClick={goBack}>{'<'} Back</a>
      </div>
      <h3>Allowable expenses:</h3>
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
        <br />
        <FontAwesomeIcon aria-hidden className="icon" icon={faBan} />Claimable expenses include: Non-claimable expenses
        include:
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
    </div>
  );
}

export default AllowableExpenses;
