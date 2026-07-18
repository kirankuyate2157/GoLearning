# CMMS and Industrial Maintenance Competitor Catalogue

**Purpose:** Feature benchmark for roadmap decisions. This is not a claim that every feature is available in every commercial plan; validate plan limits, pricing and regional support during procurement or competitive sales.

## 1. How to read this catalogue

Competitors address four levels of maturity:

```text
Level 1: work tracking      → request, work order, technician task
Level 2: maintenance system → PM, procedures, assets, inventory, reporting
Level 3: reliability system → condition, downtime, root cause, planning
Level 4: enterprise system  → EAM, ERP integration, IoT, predictive intelligence
```

The proposed product should win Levels 1–2 for a defined Indian manufacturing segment before attempting Levels 3–4.

## 2. Competitor profiles

### 2.1 MaintainX

**Positioning:** cloud CMMS for maintenance teams using web and mobile.  
**Strongest patterns to learn:** mobile execution, offline work, QR/barcode use, workflows tied to procedures and parts, and a broad integration ecosystem.

| Area | Publicly documented capability | Product lesson |
|---|---|---|
| Work | Create, assign, monitor and report work orders and work requests | The request-to-work-order conversion must be fast and obvious |
| PM | Recurring preventive work and meter/condition-triggered work | Calendar PM alone is insufficient; support meter/rule triggers later |
| Assets | Asset/location context, downtime types, criticality and history | Asset criticality and downtime reason are high-value early data fields |
| Parts | Inventory, low-stock alerts and cycle counts | Stock movements must be auditable, not only a current quantity |
| Procedures | Digital checklists and automated corrective actions | Failed inspection answers should create action, not disappear into a form |
| Mobile | QR/barcode scan, camera, timer, notifications | Technicians need scan-first and low-typing interaction |
| Offline | Cached work-order execution and later synchronization | Offline conflict/sync must be designed, not bolted on |
| Vendors | External work orders and external contacts | Vendor work needs restricted scope and sign-off |
| Enterprise | SSO, integrations, reporting and AI assistance | Enterprise readiness is a maturity path, not an MVP feature |

Sources: [MaintainX overview](https://help.getmaintainx.com/), [work orders](https://help.getmaintainx.com/about-work-orders), [mobile/web](https://help.getmaintainx.com/mobile-vs-web-overview), [offline](https://help.getmaintainx.com/offline-mode).

### 2.2 Limble

**Positioning:** intuitive CMMS focused on proactive maintenance and maintenance-team usability.  
**Strongest patterns to learn:** simplicity, work-order management, PM, asset management, inventory, resource planning and technician mobile experience.

| Area | Product lesson |
|---|---|
| Work order UX | A maintenance manager must create/assign/track work without configuration expertise |
| PM | Templates should include task steps, parts, labour and safety instruction |
| Inventory | The part must connect to asset and work order, not live in a separate spreadsheet |
| Analytics | Offer actionable metrics first: overdue work, PM completion, downtime, backlog |
| Mobile | Technician task detail must remain usable in plant conditions |

Source: [Limble product catalogue](https://limble.com/products).

### 2.3 Fiix

**Positioning:** CMMS with asset, work, parts and reliability intelligence.  
**Strongest patterns to learn:** structured asset hierarchy, part consumption/cost attribution, forecasting and maintenance data used for decisions.

| Area | Product lesson |
|---|---|
| Asset hierarchy | Model site/line/asset/component relationships clearly enough for reporting |
| Parts against work | A completed repair should record what part was consumed and its cost |
| Forecasting | Do not build forecasting until stock movement and work-order data are reliable |
| Mobile stock count | Physical counts and adjustments are necessary to trust inventory data |

Sources: [Fiix](https://fiixsoftware.com/), [feature list](https://fiixsoftware.com/wp-content/uploads/2020/05/FiixCMMSFeatureList_Final_May14.pdf).

### 2.4 UpKeep and eMaint/Fluke

**Positioning:** mobile maintenance and reliability/condition-management options across field, facility and industrial teams.  
**Lesson:** the market expects maintenance data to move from field/mobile execution into a manager’s dashboard, and eventually into condition/reliability work.

### 2.5 IBM Maximo

**Positioning:** enterprise asset-management suite for complex asset portfolios, maintenance, inspections, reliability, asset performance and investment planning.  
**Lesson:** Maximo is a horizon reference, not an MVP template. Its scope validates the long-term value of asset history + IoT + reliability data, but its implementation and enterprise complexity are not appropriate for a 20–200 asset pilot plant.

Source: [IBM Maximo Application Suite](https://www.ibm.com/products/maximo).

### 2.6 SAP EAM and Oracle EAM

**Positioning:** enterprise asset maintenance embedded in ERP-led processes.  
**Lesson:** large accounts will eventually ask for finance, purchasing, cost-centre, procurement and master-data integrations. Do not replace ERP in the initial product; provide exports and targeted integration adapters.

## 3. Detailed feature benchmark

Legend: `T` = table stake; `D` = potential differentiator; `L` = later/enterprise maturity.

| Feature family | Mature competitor expectation | Our MVP | Our differentiation | Later maturity |
|---|---|---|---|---|
| Tenant/account | Separate organization, roles, audit | T | India-region option and hands-on onboarding | SSO/SCIM |
| Assets | Hierarchy, history, documents, QR | T | Industry-specific machine templates | Component digital twin |
| Requests | Portal/QR/mobile request | T | WhatsApp-assisted intake where appropriate | AI classification |
| Work orders | Assign, priority, SLA, status, evidence | T | Minimal technician workflow | Capacity optimization |
| PM | Calendar, meter, template, escalation | T | Vertical PM libraries | Condition-based plans |
| Procedures | Digital checklists and corrective actions | T | Multilingual/local safety templates | Dynamic risk-based procedures |
| Spares | Stock, issue, vendor, reorder, count | T | GST/procurement-friendly fields | Demand forecast |
| Vendors | Contacts, service history, external work | T | Local service-partner workflows | Vendor scorecards/contracts |
| Mobile | Native/PWA, scan, camera, notifications | T | Low-bandwidth UX | Wearable/voice assistance |
| Offline | Cached work and sync | T for industrial target | Explicit conflict resolution | Full offline stores workflow |
| Reporting | PM, backlog, work, cost, downtime | T | Owner-friendly daily/weekly view | Portfolio benchmarking |
| Telemetry | Sensor/meter data and alert rules | D | MQTT/edge gateway with asset context | Condition analytics |
| AI | Summaries/recommendations | L | Explainable, local data context | Predictive model lifecycle |
| Integrations | ERP, BI, SSO, IoT | L initially | Fast spreadsheet and MQTT integration | Certified ERP connectors |

## 4. What not to copy

1. **Complex configuration before first value.** A pilot customer must import assets and create useful work within days.
2. **Feature count as positioning.** Customers buy a solved workflow, not an endless menu.
3. **AI as a substitute for data quality.** Rules and reporting beat an unsupported prediction.
4. **Desktop-only workflow.** Actual maintenance happens on the factory floor.
5. **Unbounded customization.** Allow fields/templates within clear tenant limits; avoid creating a custom product for each customer.
6. **ERP replacement ambition.** Integrate or export instead of absorbing finance/procurement scope too early.

## 5. Differentiation scorecard for discovery interviews

Ask every prospective customer to score 1–5:

| Question | Why it matters |
|---|---|
| Would you pay for assisted asset/spares setup? | Validates implementation-led motion |
| Do technicians need local language or voice/photo input? | Validates adoption wedge |
| Is WhatsApp currently used to report failures? | Validates channel integration |
| Do you need offline work at specific sites/areas? | Determines native/offline priority |
| Which ERP/accounting system is used? | Determines integration order |
| Can you connect an MQTT/PLC gateway safely? | Determines IIoT readiness |
| Which audit/customer compliance reports are painful? | Identifies vertical template opportunity |
| What financial impact did the latest major breakdown create? | Builds willingness-to-pay case |

## 6. Competitive positioning statement

> For mid-sized Indian manufacturers who currently rely on paper, Excel and WhatsApp to maintain production equipment, `[Working name]` is the maintenance-operations platform that makes preventive work, breakdown response, spares and evidence simple for technicians and visible for plant leaders. Unlike generic global CMMS products, it combines assisted local rollout, plant-ready mobile workflows and connected-machine alerting without enterprise EAM complexity.

