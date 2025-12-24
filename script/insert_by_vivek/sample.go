package main

// ============================================
// Industries
// ============================================

var Industries = []string{
	// Technology
	"software development", "cloud computing", "artificial intelligence", "machine learning",
	"cybersecurity", "data analytics", "internet of things", "blockchain", "fintech",
	"edtech", "healthtech", "e-commerce", "saas", "paas", "iaas", "devops",

	// Traditional
	"manufacturing", "retail", "healthcare", "pharmaceuticals", "biotechnology",
	"financial services", "banking", "insurance", "real estate", "construction",

	// Services
	"consulting", "legal services", "accounting", "marketing", "advertising",
	"public relations", "human resources", "staffing", "recruiting",

	// Media & Entertainment
	"media", "entertainment", "gaming", "music", "film production", "publishing",
	"broadcasting", "streaming services", "social media", "content creation",

	// Energy & Utilities
	"oil & gas", "renewable energy", "solar energy", "wind energy", "utilities",
	"electric vehicles", "clean technology", "energy storage",

	// Transportation & Logistics
	"transportation", "logistics", "supply chain", "shipping", "aviation",
	"automotive", "railways", "warehousing", "last mile delivery",

	// Consumer Goods
	"consumer electronics", "food & beverage", "fashion", "apparel", "cosmetics",
	"home goods", "furniture", "luxury goods", "sports equipment",

	// Other
	"telecommunications", "aerospace", "defense", "agriculture", "mining",
	"hospitality", "travel", "tourism", "education", "non-profit",
}

// ============================================
// Keywords
// ============================================

var Keywords = []string{
	// Business Models
	"b2b", "b2c", "d2c", "marketplace", "platform", "subscription", "freemium",
	"enterprise", "smb", "startup", "scale-up", "unicorn", "decacorn",

	// Growth Stage
	"pre-seed", "seed", "series a", "series b", "series c", "late stage",
	"ipo", "public", "private", "bootstrapped", "venture-backed",

	// Technology Focus
	"api-first", "mobile-first", "cloud-native", "open source", "proprietary",
	"low-code", "no-code", "ai-powered", "data-driven", "automation",

	// Market Position
	"market leader", "disruptor", "challenger", "emerging", "established",
	"global", "regional", "local", "niche", "mass market",

	// Business Attributes
	"fast-growing", "profitable", "high-growth", "sustainable", "innovation",
	"remote-first", "hybrid", "distributed", "agile", "lean",

	// Industry Specific
	"regtech", "insurtech", "proptech", "agritech", "foodtech", "cleantech",
	"biotech", "medtech", "legaltech", "hrtech", "martech", "adtech",
}

// ============================================
// Technologies
// ============================================

var Technologies = []string{
	// Programming Languages
	"python", "javascript", "typescript", "java", "go", "rust", "c++", "c#",
	"ruby", "php", "kotlin", "swift", "scala", "r", "matlab",

	// Frontend
	"react", "vue.js", "angular", "next.js", "nuxt.js", "svelte", "jquery",
	"tailwind css", "bootstrap", "material ui", "chakra ui",

	// Backend
	"node.js", "django", "flask", "fastapi", "spring boot", "express.js",
	"ruby on rails", "laravel", "asp.net", "gin", "echo",

	// Databases
	"postgresql", "mysql", "mongodb", "redis", "elasticsearch", "cassandra",
	"dynamodb", "firebase", "supabase", "cockroachdb", "timescaledb",

	// Cloud & Infrastructure
	"aws", "google cloud", "azure", "digitalocean", "heroku", "vercel",
	"netlify", "cloudflare", "kubernetes", "docker", "terraform",

	// Data & ML
	"tensorflow", "pytorch", "scikit-learn", "pandas", "numpy", "spark",
	"hadoop", "kafka", "airflow", "dbt", "snowflake", "databricks",

	// DevOps & Tools
	"git", "github", "gitlab", "jenkins", "circleci", "github actions",
	"ansible", "prometheus", "grafana", "datadog", "new relic", "sentry",

	// Communication & Collaboration
	"slack", "microsoft teams", "zoom", "notion", "confluence", "jira",
	"linear", "asana", "monday.com", "figma", "miro",
}

// ============================================
// Company Name Components
// ============================================

var CompanyNamePrefixes = []string{
	"Tech", "Digital", "Smart", "Cloud", "Data", "Cyber", "Net", "Web",
	"App", "Soft", "Info", "Inno", "Pro", "Prime", "Elite", "Alpha",
	"Beta", "Gamma", "Delta", "Omega", "Hyper", "Ultra", "Meta", "Neo",
	"Quantum", "Vertex", "Apex", "Peak", "Summit", "Zenith", "Stellar",
	"Astro", "Cosmo", "Global", "Inter", "Trans", "Multi", "Omni", "Uni",
	"Flex", "Agile", "Swift", "Rapid", "Turbo", "Power", "Force", "Core",
	"Fusion", "Synergy", "Vector", "Matrix", "Nexus", "Pulse", "Wave",
}

var CompanyNameSuffixes = []string{
	"Labs", "Works", "Systems", "Solutions", "Technologies", "Tech", "Soft",
	"Ware", "Net", "Cloud", "Data", "Logic", "Minds", "Dynamics", "Ventures",
	"Group", "Corp", "Inc", "Co", "Hub", "Space", "Studio", "Digital",
	"Analytics", "AI", "Intelligence", "Innovations", "Partners", "Global",
	"Interactive", "Media", "Networks", "Services", "Consulting", "Platform",
	"Apps", "Mobile", "Connect", "Link", "Bridge", "Gate", "Port", "Base",
}

var CompanyNameWords = []string{
	"Acme", "Apex", "Horizon", "Pinnacle", "Summit", "Vanguard", "Pioneer",
	"Catalyst", "Momentum", "Velocity", "Altitude", "Latitude", "Longitude",
	"Spectrum", "Prism", "Helix", "Spiral", "Orbit", "Galaxy", "Nova",
	"Phoenix", "Titan", "Atlas", "Zeus", "Apollo", "Athena", "Hermes",
	"Artemis", "Poseidon", "Hera", "Mercury", "Venus", "Mars", "Jupiter",
	"Saturn", "Neptune", "Pluto", "Orion", "Andromeda", "Sirius", "Polaris",
	"Aurora", "Eclipse", "Equinox", "Solstice", "Ember", "Blaze", "Spark",
	"Thunder", "Lightning", "Storm", "Tempest", "Cyclone", "Tornado", "Vortex",
}

// ============================================
// Location Data
// ============================================

var Countries = []string{
	"usa", "uk", "canada", "germany", "france", "india", "australia",
	"japan", "brazil", "mexico", "netherlands", "sweden", "singapore",
	"israel", "south korea", "china", "spain", "italy", "switzerland",
}

var StatesByCountry = map[string][]string{
	"usa":         {"ca", "ny", "tx", "fl", "il", "pa", "oh", "ga", "nc", "mi", "wa", "ma", "co", "az", "va"},
	"uk":          {"england", "scotland", "wales", "northern ireland"},
	"canada":      {"on", "bc", "qc", "ab", "mb", "sk", "ns", "nb"},
	"germany":     {"bavaria", "berlin", "hamburg", "hesse", "north rhine-westphalia", "baden-württemberg"},
	"france":      {"île-de-france", "provence-alpes-côte d'azur", "auvergne-rhône-alpes", "occitanie"},
	"india":       {"maharashtra", "karnataka", "delhi", "tamil nadu", "telangana", "gujarat", "west bengal"},
	"australia":   {"nsw", "vic", "qld", "wa", "sa", "tas"},
	"japan":       {"tokyo", "osaka", "kanagawa", "aichi", "fukuoka", "hokkaido"},
	"brazil":      {"são paulo", "rio de janeiro", "minas gerais", "bahia", "rio grande do sul"},
	"mexico":      {"cdmx", "jalisco", "nuevo león", "estado de méxico", "puebla"},
	"netherlands": {"north holland", "south holland", "utrecht", "north brabant"},
	"sweden":      {"stockholm", "västra götaland", "skåne", "uppsala"},
	"singapore":   {"central region", "east region", "north region", "west region"},
	"israel":      {"tel aviv", "jerusalem", "haifa", "central"},
	"south korea": {"seoul", "gyeonggi", "busan", "incheon"},
	"china":       {"beijing", "shanghai", "guangdong", "zhejiang", "jiangsu"},
	"spain":       {"madrid", "catalonia", "andalusia", "valencia"},
	"italy":       {"lombardy", "lazio", "veneto", "emilia-romagna", "piedmont"},
	"switzerland": {"zürich", "geneva", "bern", "basel", "vaud"},
}

var CitiesByCountry = map[string][]string{
	"usa":         {"new york", "san francisco", "los angeles", "chicago", "boston", "seattle", "austin", "denver", "miami", "atlanta"},
	"uk":          {"london", "manchester", "birmingham", "edinburgh", "bristol", "leeds", "glasgow", "cambridge", "oxford"},
	"canada":      {"toronto", "vancouver", "montreal", "calgary", "ottawa", "edmonton", "waterloo"},
	"germany":     {"berlin", "munich", "hamburg", "frankfurt", "cologne", "stuttgart", "düsseldorf"},
	"france":      {"paris", "lyon", "marseille", "toulouse", "nice", "nantes", "bordeaux"},
	"india":       {"bangalore", "mumbai", "delhi", "hyderabad", "chennai", "pune", "gurgaon", "noida"},
	"australia":   {"sydney", "melbourne", "brisbane", "perth", "adelaide", "canberra"},
	"japan":       {"tokyo", "osaka", "yokohama", "nagoya", "fukuoka", "kyoto", "sapporo"},
	"brazil":      {"são paulo", "rio de janeiro", "belo horizonte", "brasília", "curitiba", "porto alegre"},
	"mexico":      {"mexico city", "guadalajara", "monterrey", "puebla", "tijuana", "león"},
	"netherlands": {"amsterdam", "rotterdam", "the hague", "utrecht", "eindhoven"},
	"sweden":      {"stockholm", "gothenburg", "malmö", "uppsala"},
	"singapore":   {"singapore"},
	"israel":      {"tel aviv", "jerusalem", "haifa", "herzliya", "ra'anana"},
	"south korea": {"seoul", "busan", "incheon", "daegu", "daejeon"},
	"china":       {"beijing", "shanghai", "shenzhen", "hangzhou", "guangzhou", "chengdu"},
	"spain":       {"madrid", "barcelona", "valencia", "seville", "bilbao"},
	"italy":       {"milan", "rome", "turin", "florence", "bologna", "naples"},
	"switzerland": {"zürich", "geneva", "basel", "bern", "lausanne"},
}

var DomainExtensions = []string{
	".com", ".io", ".co", ".ai", ".tech", ".app", ".dev", ".cloud",
	".net", ".org", ".biz", ".info", ".solutions", ".digital", ".software",
}

// ============================================
// Company Size and Funding Data
// ============================================

type RangeInt64 struct {
	Min int64
	Max int64
}

var EmployeeCountRanges = map[string]RangeInt64{
	"startup":    {1, 50},
	"small":      {51, 200},
	"medium":     {201, 1000},
	"large":      {1001, 10000},
	"enterprise": {10001, 500000},
}

var AnnualRevenueRanges = map[string]RangeInt64{
	"startup":    {0, 5_000_000},
	"small":      {5_000_001, 50_000_000},
	"medium":     {50_000_001, 500_000_000},
	"large":      {500_000_001, 5_000_000_000},
	"enterprise": {5_000_000_001, 100_000_000_000},
}

var TotalFundingRanges = map[string]RangeInt64{
	"bootstrapped": {0, 100_000},
	"pre_seed":     {100_001, 500_000},
	"seed":         {500_001, 3_000_000},
	"series_a":     {3_000_001, 15_000_000},
	"series_b":     {15_000_001, 50_000_000},
	"series_c":     {50_000_001, 150_000_000},
	"late_stage":   {150_000_001, 1_000_000_000},
	"public":       {0, 0},
}

var CompanySizeWeights = map[string]float64{
	"startup":    0.30,
	"small":      0.30,
	"medium":     0.20,
	"large":      0.12,
	"enterprise": 0.08,
}

var FundingStageWeights = map[string]float64{
	"bootstrapped": 0.20,
	"pre_seed":     0.10,
	"seed":         0.20,
	"series_a":     0.18,
	"series_b":     0.12,
	"series_c":     0.08,
	"late_stage":   0.07,
	"public":       0.05,
}

var FundingStages = []string{
	"Bootstrapped", "Pre-Seed", "Seed", "Series A", "Series B",
	"Series C", "Series D", "Series E", "Late Stage", "IPO", "Public",
}

// ============================================
// Contact Data
// ============================================

var Titles = []string{
	// Engineering & Tech
	"Software Engineer", "Senior Software Engineer", "Staff Engineer", "Principal Engineer",
	"Full Stack Developer", "Frontend Developer", "Backend Developer", "Mobile Developer",
	"DevOps Engineer", "Site Reliability Engineer", "Cloud Architect", "Infrastructure Engineer",
	"Security Engineer", "Cybersecurity Specialist", "Penetration Tester", "Security Architect",
	"Machine Learning Engineer", "AI Engineer", "Data Engineer", "Data Scientist", "MLOps Engineer",
	"QA Engineer", "Test Automation Engineer", "Quality Assurance Lead", "Performance Engineer",
	"Blockchain Developer", "Smart Contract Developer", "Web3 Engineer", "Cryptocurrency Analyst",

	// Product & Design
	"Product Manager", "Senior Product Manager", "Product Owner", "Technical Product Manager",
	"UI/UX Designer", "Product Designer", "Visual Designer", "Interaction Designer",
	"Design Systems Lead", "Creative Director", "User Experience Researcher", "Design Strategist",

	// Data & Analytics
	"Data Analyst", "Business Intelligence Analyst", "Data Architect", "Analytics Engineer",
	"Business Analyst", "Product Analyst", "Growth Analyst", "Marketing Analyst",

	// Marketing & Growth
	"Marketing Manager", "Digital Marketing Specialist", "Content Marketing Manager",
	"Growth Marketing Manager", "SEO Specialist", "Social Media Manager", "Brand Manager",
	"Marketing Automation Specialist", "Performance Marketing Manager", "Community Manager",
	"Influencer Relations Manager", "Content Creator", "Copywriter", "Marketing Technologist",

	// Sales & Business Development
	"Sales Representative", "Account Executive", "Sales Development Representative",
	"Business Development Manager", "Enterprise Sales Manager", "Sales Engineer",
	"Customer Success Manager", "Account Manager", "Partnership Manager", "Channel Manager",

	// Operations & Strategy
	"Operations Manager", "Business Operations Analyst", "Strategy Manager",
	"Program Manager", "Project Manager", "Scrum Master", "Agile Coach",
	"Process Improvement Specialist", "Operations Analyst", "Supply Chain Manager",

	// HR & People
	"HR Manager", "Talent Acquisition Specialist", "People Operations Manager",
	"HR Business Partner", "Learning & Development Manager", "Compensation Analyst",
	"Employee Experience Manager", "Diversity & Inclusion Manager", "Recruiter",

	// Finance & Accounting
	"Finance Analyst", "Financial Controller", "Accounting Manager", "FP&A Analyst",
	"Treasury Analyst", "Tax Specialist", "Auditor", "Financial Planning Manager",

	// Executive & Leadership
	"CEO", "CTO", "CFO", "COO", "CMO", "CPO", "CHRO", "Chief Data Officer",
	"VP of Engineering", "VP of Product", "VP of Sales", "VP of Marketing",
	"Director of Engineering", "Director of Product", "Director of Sales",

	// Specialized & Emerging Roles
	"Customer Experience Designer", "Voice of Customer Analyst", "Technical Writer",
	"Developer Advocate", "Developer Relations Manager", "Open Source Program Manager",
	"Ethics Engineer", "Privacy Engineer", "Compliance Manager", "Legal Counsel",
	"Research Scientist", "Applied Scientist", "Technical Program Manager",
	"Solutions Architect", "Customer Success Engineer", "Implementation Specialist",
}

var Departments = []string{
	"engineering", "product", "sales", "marketing", "hr", "finance",
	"operations", "customer success", "support", "legal",
}

var SeniorityLevels = []string{
	"junior", "mid", "senior", "lead", "principal", "executive",
}

var EmailStatuses = []string{
	"verified", "unverified", "bounced", "invalid",
}

var ContactStages = []string{
	"lead", "contacted", "qualified", "proposal", "negotiation", "closed won", "closed lost",
}