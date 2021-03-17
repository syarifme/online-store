# online-store

# the cause of the bad rating
The inventory quantities are ofter misreported, this causes the number of orders to be greater than the number of stock items. The intentory application can't compare the number of order and the number of quantity, so it can't limit the order number. After reaching the limit order number or order number equals to number of stock items, the inventory application does not cancel orders.

# solution proposed
To prevent future incidents, the inventory application have to hold the number of stock items before flash sale begins. When an order comes with a certain quantity, the number of stock items is reduced according to the order quantity requested. If the order quantity is greater than the number of stock items or zero number of stock items, it will return an error and the order fails.