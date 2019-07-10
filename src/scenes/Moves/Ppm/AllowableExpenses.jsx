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
      <h3 className="title">Storage & Moving Expenses</h3>
      <p>
        <strong>Storage expenses</strong> are a special expense that is <strong>reimbursable</strong> for up to 90 days.
        You can be directly repaid for those expenses.
      </p>
      <p>
        The IRS considers the rest of your PPM payment as taxable income. You’ll receive a separate W-2 for any PPM
        payment.
      </p>
      <p>
        <strong>Moving-related expenses</strong> can be <strong>claimed</strong> in order to reduce the taxable amount
        of your payment. Your{' '}
        <a href="https://installations.militaryonesource.mil/search" target="_blank">
          local finance office
        </a>{' '}
        or a tax professional can help you identify qualifying expenses. You can also consult{' '}
        <a href="https://www.irs.gov/publications/p521" target="_blank">
          IRS Publication 521
        </a>{' '}
        for authoritative information.
      </p>
      <p>
        <strong>Save your receipts.</strong> It’s better to have receipts you don’t need than to need receipts you don’t
        have.
      </p>
      <hr className="divider" />
      <div className="bullet-li-header">
        <FontAwesomeIcon aria-hidden className="icon" icon={faCheck} />Some commonly claimed moving expenses:
      </div>
      <ul>
        <li>Consumable packing materials</li>
        <li>Contracted expenses</li>
        <li>Oil</li>
        <li>Rental equipment</li>
        <li>Tolls</li>
        <li>Weighing fees</li>
        <li>Gas, exceeding travel allowance (see below)</li>
      </ul>
      <div className="divider dashed-divider" />
      <div className="bullet-li-header">
        <FontAwesomeIcon aria-hidden className="icon" icon={faBan} />Some common expenses that are <em>not</em>{' '}
        claimable or reimbursable:
      </div>
      <ul>
        <li>Animal costs (kennels, transportation)</li>
        <li>Extra drivers</li>
        <li>Hitch fees and tow bars</li>
        <li>Locks</li>
        <li>Meals and lodging</li>
        <li>Moving insurance</li>
        <li>Oil change and routine maintenance</li>
        <li>Purchased auto transporters and dollies</li>
        <li>Sales tax</li>
        <li>Tire chains</li>
        <li>Gas, under travel allowance (see details following)</li>
      </ul>
      <hr className="divider" />
      <section style={{ marginBottom: '1.5em' }}>
        <strong>Gas and fuel expenses</strong>
        <p>
          Fuel expenses may not be claimed for a PPM unless they exceed the amount paid for mileage and per diem fees on
          your travel pay. The IRS does not allow you to claim an expense if you were already paid for the mileage.
          Doing so could result in an IRS audit.
        </p>
        <p>
          If your fuel costs exceed the per diem payment you received, however, you may claim the portion that exceeds
          that amount.
        </p>
      </section>
      <section style={{ marginBottom: '3.5em' }}>
        <strong>When are receipts required?</strong>
        <p>You must have receipts for contracted expenses, storage facilities, and any expense over $75.</p>
        <p>
          You will need receipts for expenses under $75 if multiple expenses in that same category add up to more than
          $75.
        </p>
        <p>
          Again, it’s better to have receipts you don’t need than to be missing receipts that you do need. We recommend{' '}
          <strong>saving all your moving receipts.</strong>
        </p>
        <p>
          If you are missing a receipt, you can go online and print a new copy of your receipt (if you can). Otherwise,
          write and sign a statement that explains why the receipt is missing. Contact your{' '}
          <a href="https://installations.militaryonesource.mil/search" target="_blank">
            local finance office
          </a>{' '}
          for assistance.
        </p>
      </section>
      <div className="usa-grid button-bar">
        <button className="usa-button-secondary" onClick={goBack}>
          Back
        </button>
      </div>
    </div>
  );
}

export default AllowableExpenses;
