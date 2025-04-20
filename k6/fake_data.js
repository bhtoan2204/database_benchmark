import { check } from "k6";
import { randomIntBetween } from "https://jslib.k6.io/k6-utils/1.2.0/index.js";
import http from "k6/http";
import faker from "https://cdnjs.cloudflare.com/ajax/libs/Faker/3.1.0/faker.min.js";

const API_URL = "http://localhost:8080/api/v1";
const NUM_PRODUCTS = 1000000; // 1 triệu products
const BATCH_SIZE = 100; // Giảm xuống 100 products mỗi batch

const categories = [
  { id: "electronics", name: "Electronics" },
  { id: "computers", name: "Computers & Tablets" },
  { id: "mobile", name: "Mobile Phones" },
  { id: "cameras", name: "Cameras & Photo" },
  { id: "audio", name: "Audio & Headphones" },
  { id: "gaming", name: "Gaming" },
  { id: "smart-home", name: "Smart Home" },
  { id: "wearables", name: "Wearables" },
  { id: "tv-video", name: "TV & Video" },
  { id: "accessories", name: "Accessories" },
  { id: "health", name: "Health & Beauty" },
  { id: "home", name: "Home & Garden" },
  { id: "toys", name: "Toys & Games" },
  { id: "sports", name: "Sports & Outdoors" },
  { id: "automotive", name: "Automotive" },
  { id: "groceries", name: "Groceries" },
  { id: "books", name: "Books" },
  { id: "music", name: "Music" },
  { id: "office", name: "Office Supplies" },
  { id: "art", name: "Art & Crafts" },
  { id: "pet", name: "Pet Supplies" },
  { id: "baby", name: "Baby & Toddler" },
  { id: "jewelry", name: "Jewelry & Watches" },
  { id: "tools", name: "Tools & Home Improvement" },
  { id: "garden", name: "Garden & Outdoor" },
];

const brands = [
  { id: "apple", name: "Apple" },
  { id: "samsung", name: "Samsung" },
  { id: "google", name: "Google" },
  { id: "microsoft", name: "Microsoft" },
  { id: "sony", name: "Sony" },
  { id: "lg", name: "LG" },
  { id: "huawei", name: "Huawei" },
  { id: "xiaomi", name: "Xiaomi" },
  { id: "oneplus", name: "OnePlus" },
  { id: "nokia", name: "Nokia" },
  { id: "asus", name: "Asus" },
  { id: "dell", name: "Dell" },
  { id: "hp", name: "HP" },
  { id: "lenovo", name: "Lenovo" },
  { id: "acer", name: "Acer" },
  { id: "toshiba", name: "Toshiba" },
];

const productTypes = {
  mobile: ["smartphone", "tablet", "smartwatch", "earphone", "speaker"],
  computer: ["laptop", "desktop", "monitor", "keyboard", "mouse"],
  camera: ["camera", "lens", "tripod", "camera bag", "camera accessory"],
  headphone: ["earphone", "headphone", "speaker", "earphone accessory"],
  speaker: ["speaker", "speaker accessory"],
  game: ["game console", "game console accessory", "game console game"],
  smartHome: ["smart home device", "smart home accessory"],
  wearables: ["smartwatch", "fitness band", "smartwatch accessory"],
  tvVideo: ["tv", "monitor", "speaker", "tv accessory"],
  accessories: ["accessory", "accessory accessory"],
  health: ["health device", "health accessory"],
};

const tags = [
  "new-arrival",
  "bestseller",
  "featured",
  "sale",
  "trending",
  "premium",
  "limited-edition",
  "exclusive",
  "eco-friendly",
  "wireless",
  "bluetooth",
  "gaming",
  "professional",
  "budget-friendly",
  "luxury",
  "refurbished",
  "open-box",
  "clearance",
  "pre-order",
  "bundle",
];

function generateProductName(brand, type, category) {
  const year = new Date().getFullYear();
  const modelNumber = faker.random.alphaNumeric(4).toUpperCase();
  const series = faker.random.arrayElement([
    "Pro",
    "Elite",
    "Plus",
    "Max",
    "Ultra",
    "Lite",
    "",
  ]);

  switch (category.id) {
    case "mobile":
      return `${brand.name} ${type} ${series} ${modelNumber} (${year})`;
    case "computers":
      return `${
        brand.name
      } ${type} ${series} ${modelNumber} ${faker.random.arrayElement([
        "i5",
        "i7",
        "i9",
        "Ryzen 5",
        "Ryzen 7",
        "Ryzen 9",
      ])}`;
    case "gaming":
      return `${brand.name} ${type} ${series} RGB ${modelNumber}`;
    default:
      return `${brand.name} ${type} ${series} ${modelNumber}`;
  }
}

function generateDescription(product, brand, type, category) {
  const features = [];

  // Add category-specific features
  switch (category.id) {
    case "mobile":
      features.push(
        faker.random.arrayElement(["5G", "4G LTE"]),
        `${randomIntBetween(4, 12)}GB RAM`,
        `${randomIntBetween(64, 512)}GB Storage`,
        `${randomIntBetween(4000, 5000)}mAh Battery`,
        faker.random.arrayElement(["AMOLED", "Retina", "LCD"]) + " Display"
      );
      break;
    case "computers":
      features.push(
        `${randomIntBetween(8, 64)}GB RAM`,
        `${randomIntBetween(256, 2048)}GB ${faker.random.arrayElement([
          "SSD",
          "NVMe SSD",
        ])}`,
        faker.random.arrayElement(["Windows 11", "macOS", "Chrome OS"]),
        `${randomIntBetween(13, 17)}" Display`
      );
      break;
    case "audio":
      features.push(
        "Active Noise Cancellation",
        `${randomIntBetween(20, 40)} Hours Battery Life`,
        faker.random.arrayElement(["Bluetooth 5.0", "Bluetooth 5.2"]),
        "Touch Controls"
      );
      break;
  }

  const description =
    `Experience the next level of ${category.name.toLowerCase()} with the ${
      product.name
    }. ` +
    `This ${
      brand.tier
    }-tier ${type.toLowerCase()} combines cutting-edge technology with ${
      brand.name
    }'s signature quality. ` +
    `Key features include: ${features.join(", ")}. ` +
    faker.lorem.paragraph();

  return description;
}

function generateAttributes(category, type) {
  const commonAttrs = {
    color: faker.random.arrayElement([
      "Black",
      "White",
      "Silver",
      "Space Gray",
      "Midnight Blue",
      "Rose Gold",
    ]),
    model_year: randomIntBetween(2022, 2024),
    warranty: faker.random.arrayElement(["1 Year", "2 Years", "3 Years"]),
  };

  switch (category.id) {
    case "mobile":
      return {
        ...commonAttrs,
        screen_size: `${(randomIntBetween(60, 70) / 10).toFixed(1)}"`,
        ram: `${randomIntBetween(4, 12)}GB`,
        storage: `${randomIntBetween(64, 512)}GB`,
        camera: `${randomIntBetween(12, 108)}MP`,
        battery: `${randomIntBetween(4000, 5000)}mAh`,
        os_version: faker.random.arrayElement([
          "iOS 17",
          "Android 14",
          "Android 13",
        ]),
      };
    case "computers":
      return {
        ...commonAttrs,
        processor: faker.random.arrayElement([
          "Intel Core i5",
          "Intel Core i7",
          "Intel Core i9",
          "AMD Ryzen 5",
          "AMD Ryzen 7",
          "AMD Ryzen 9",
        ]),
        ram: `${randomIntBetween(8, 64)}GB`,
        storage: `${randomIntBetween(256, 2048)}GB`,
        screen_size: `${randomIntBetween(13, 32)}"`,
        graphics: faker.random.arrayElement([
          "NVIDIA RTX 3060",
          "NVIDIA RTX 3070",
          "NVIDIA RTX 4060",
          "AMD Radeon RX 6600",
          "Intel Iris Xe",
          "Apple M2",
        ]),
      };
    case "audio":
      return {
        ...commonAttrs,
        driver_size: `${randomIntBetween(30, 50)}mm`,
        frequency_response: `20Hz-20kHz`,
        battery_life: `${randomIntBetween(20, 40)} hours`,
        bluetooth_version: faker.random.arrayElement(["5.0", "5.2", "5.3"]),
      };
    default:
      return commonAttrs;
  }
}

function generatePrice(brand, category) {
  const tierMultiplier = {
    premium: 2,
    mid: 1.5,
    budget: 1,
    gaming: 1.8,
  };

  const basePrice = {
    mobile: { min: 19900, max: 129900 },
    computers: { min: 49900, max: 299900 },
    audio: { min: 9900, max: 39900 },
    gaming: { min: 29900, max: 199900 },
    cameras: { min: 39900, max: 299900 },
    accessories: { min: 1990, max: 19900 },
  };

  const priceRange = basePrice[category.id] || basePrice.accessories;
  const multiplier = tierMultiplier[brand.tier] || 1;

  return Math.floor(
    randomIntBetween(priceRange.min, priceRange.max) * multiplier
  );
}

function generateSKU(index, brand, category) {
  const brandPrefix = brand.name.substring(0, 3).toUpperCase();
  const catPrefix = category.id.substring(0, 3).toUpperCase();
  const serialNum = index.toString().padStart(6, "0");
  return `${brandPrefix}${catPrefix}${serialNum}`;
}

function getRandomElements(array, n) {
  const shuffled = array.sort(() => 0.5 - Math.random());
  return shuffled.slice(0, n);
}

function generateRandomProduct(index) {
  // Select random category and brand
  const category = faker.random.arrayElement(categories);
  const brand = faker.random.arrayElement(brands);

  // Get product type based on category
  const availableTypes = productTypes[category.id] || productTypes.accessories;
  const type = faker.random.arrayElement(availableTypes);

  // Generate product name
  const name = generateProductName(brand, type, category);

  // Select 1-3 random categories (including main category)
  const numCategories = randomIntBetween(1, 3);
  const additionalCategories = getRandomElements(
    categories.filter((c) => c.id !== category.id),
    numCategories - 1
  );
  const selectedCategories = [category, ...additionalCategories];

  // Select 2-5 random tags
  const numTags = randomIntBetween(2, 5);
  const selectedTags = getRandomElements(tags, numTags);

  // Generate random images
  const numImages = randomIntBetween(3, 6);
  const images = Array(numImages)
    .fill(0)
    .map((_, i) => {
      const width = randomIntBetween(800, 1200);
      const height = randomIntBetween(800, 1200);
      return `https://picsum.photos/id/${index * 10 + i}/${width}/${height}`;
    });

  // Generate dates
  const now = new Date();
  const createdAt = faker.date.past(2); // Within the last 2 years
  const updatedAt = faker.date.between(createdAt, now);

  // Calculate price based on brand tier and category
  const price = generatePrice(brand, category);

  // Generate stock based on price and brand tier
  const baseStock =
    brand.tier === "premium"
      ? randomIntBetween(5, 50)
      : randomIntBetween(20, 200);

  return {
    id: `prod-${faker.random.uuid()}`,
    name: name,
    sku: generateSKU(index, brand, category),
    description: generateDescription({ name }, brand, type, category),
    price: price,
    category: selectedCategories,
    tags: selectedTags,
    images: images,
    stock: baseStock,
    brand: brand,
    attributes: generateAttributes(category, type),
    rating_avg: Number((Math.random() * 2 + 3).toFixed(1)), // Random rating between 3.0 and 5.0
    rating_count: randomIntBetween(0, 2000),
    created_at: createdAt.toISOString(),
    updated_at: updatedAt.toISOString(),
    deleted_at: null,
  };
}

export const options = {
  scenarios: {
    bulk_insert: {
      executor: "shared-iterations",
      vus: 20, // Tăng số VUs vì batch size nhỏ hơn
      iterations: Math.ceil(NUM_PRODUCTS / BATCH_SIZE), // = 10000 iterations
      maxDuration: "24h",
    },
  },
  batch: {
    batchPerHost: 10, // Tăng số concurrent requests vì mỗi request nhẹ hơn
  },
  thresholds: {
    http_req_failed: ["rate<0.01"], // Dưới 1% requests thất bại
    http_req_duration: ["p(95)<3000"], // 95% requests dưới 3s
  },
};

export default function () {
  // Generate một batch 100 sản phẩm
  const currentBatch = Array(BATCH_SIZE)
    .fill(null)
    .map((_, index) => generateRandomProduct(__ITER * BATCH_SIZE + index + 1));

  // Gửi batch lên API
  const payload = JSON.stringify(currentBatch);
  const params = {
    headers: {
      "Content-Type": "application/json",
    },
    timeout: "30s", // Giảm timeout vì payload nhỏ hơn
  };

  const response = http.post(`${API_URL}/products/bulk`, payload, params);

  // Kiểm tra response
  const success = check(response, {
    "is status 201": (r) => r.status === 201,
    "has successful response": (r) => {
      try {
        const body = JSON.parse(r.body);
        return body.count === currentBatch.length;
      } catch (e) {
        return false;
      }
    },
  });

  // Log progress mỗi 10,000 records (100 batches)
  if (success) {
    const totalProcessed = (__ITER + 1) * BATCH_SIZE;
    if (totalProcessed % 10000 === 0 || totalProcessed >= NUM_PRODUCTS) {
      const elapsed = new Date() - __ENV.startTime;
      const rate = totalProcessed / (elapsed / 1000);
      console.log(
        `Processed ${Math.min(
          totalProcessed,
          NUM_PRODUCTS
        )}/${NUM_PRODUCTS} products ` +
          `(${rate.toFixed(2)} products/sec) - ` +
          `Batch ${__ITER + 1}/${Math.ceil(NUM_PRODUCTS / BATCH_SIZE)}`
      );
    }
  }

  // Thêm tiny delay giữa các requests để tránh overload
  if (__ITER % 10 === 0) {
    // Mỗi 10 batches (1000 products)
    sleep(0.1); // 100ms delay
  }
}

// Thêm tracking thời gian bắt đầu
export function setup() {
  return { startTime: new Date() };
}

export function handleSummary(data) {
  const metrics = {
    totalRequests: data.metrics.iterations.values.count,
    failedRequests: data.metrics.http_req_failed.values.passes,
    avgRequestDuration: data.metrics.http_req_duration.values.avg,
    p95RequestDuration: data.metrics.http_req_duration.values["p(95)"],
    totalDuration: (new Date() - __ENV.startTime) / 1000,
  };

  return {
    stdout: textSummary(data, { indent: " ", enableColors: true }),
    "./k6-summary.json": JSON.stringify(
      {
        ...data,
        customMetrics: metrics,
      },
      null,
      2
    ),
    "./k6-metrics.txt": `
Total Products Processed: ${metrics.totalRequests * BATCH_SIZE}
Failed Requests: ${metrics.failedRequests}
Average Request Duration: ${metrics.avgRequestDuration.toFixed(2)}ms
P95 Request Duration: ${metrics.p95RequestDuration.toFixed(2)}ms
Total Duration: ${metrics.totalDuration.toFixed(2)}s
Average Rate: ${(
      (metrics.totalRequests * BATCH_SIZE) /
      metrics.totalDuration
    ).toFixed(2)} products/sec
        `.trim(),
  };
}
