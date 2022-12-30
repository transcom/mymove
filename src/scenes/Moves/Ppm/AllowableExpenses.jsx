import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useNavigate } from 'react-router-dom';

function AllowableExpenses(props) {
  const navigate = useNavigate();
  function goBack() {
    navigate(-1);
  }

  return (
    <div className="grid-container usa-prose allowable-expenses-container">
      <div className="grid-row">
        <div className="grid-col-12">
          <div>
            <a onClick={goBack} className="usa-link">
              {'<'} Back
            </a>
          </div>
          <h1 className="title">Storage & Moving Expenses</h1>
          <p>
            <strong>Storage expenses</strong> are a special expense that is <strong>reimbursable</strong> for up to 90
            days. You can be directly repaid for those expenses.
          </p>
          <p>
            The IRS considers the rest of your PPM payment as taxable income. You’ll receive a separate W-2 for any PPM
            payment.
          </p>
          <p>
            <strong>Moving-related expenses</strong> can be <strong>claimed</strong> in order to reduce the taxable
            amount of your payment. Your{' '}
            <a
              href="https://installations.militaryonesource.mil/search"
              target="_blank"
              rel="noopener noreferrer"
              className="usa-link"
            >
              local finance office
            </a>{' '}
            or a tax professional can help you identify qualifying expenses. You can also consult{' '}
            <a
              href="https://www.irs.gov/publications/p521"
              target="_blank"
              rel="noopener noreferrer"
              className="usa-link"
            >
              IRS Publication 521
            </a>{' '}
            for authoritative information.
          </p>
          <p>
            <strong>Save your receipts.</strong> It’s better to have receipts you don’t need than to need receipts you
            don’t have.
          </p>
          <hr className="divider" />
          <div className="bullet-li-header">
            <FontAwesomeIcon aria-hidden className="icon" icon="check" />
            Some commonly claimed moving expenses:
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
            <FontAwesomeIcon aria-hidden className="icon" icon="ban" />
            Some common expenses that are <em>not</em> claimable or reimbursable:
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
            <p>Gas and fuel expenses are not reimbursable.</p>
            <p>If you rented a vehicle to perform your move, you can claim gas expenses for tax purposes.</p>
            <p>
              You can not claim expenses for fuel for your own vehicles. You will be reimbursed for that fuel via DTS
              when you claim your mileage. The IRS does not allow you to claim an expense twice, and may audit you if
              you do so.
            </p>
            <p>
              There is one rare exception: If your fuel expenses exceed the amount paid for mileage and per diem fees on
              your travel pay. You may claim the portion of your fuel expenses for your own vehicles that exceeds that
              amount.
            </p>
          </section>
          <section style={{ marginBottom: '3.5em' }}>
            <strong>When are receipts required?</strong>
            <p>You must have receipts for contracted expenses, storage facilities, and any expense over $75.</p>
            <p>
              You will need receipts for expenses under $75 if multiple expenses in that same category add up to more
              than $75.
            </p>
            <p>
              Again, it’s better to have receipts you don’t need than to be missing receipts that you do need. We
              recommend <strong>saving all your moving receipts.</strong>
            </p>
            <p>
              If you are missing a receipt, you can go online and print a new copy of your receipt (if you can).
              Otherwise, write and sign a statement that explains why the receipt is missing. Contact your{' '}
              <a
                href="https://installations.militaryonesource.mil/search"
                target="_blank"
                rel="noopener noreferrer"
                className="usa-link"
              >
                local finance office
              </a>{' '}
              for assistance.
            </p>
          </section>
          <div className="usa-grid button-bar">
            <button className="usa-button usa-button--secondary" onClick={goBack}>
              Back
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default AllowableExpenses;
